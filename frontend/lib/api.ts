// API関連のユーティリティ関数

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080";

// バックエンドのエンティティに完全に合わせた型定義
export interface User {
  id: number;
  email: string;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  name: string;
  password: string;
}

// バックエンドの AuthResponse に完全に合わせる
export interface AuthResponse {
  user: User;
  token: string;
}

export interface PaginationResponse {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

export interface UsersResponse {
  users: User[];
  pagination: PaginationResponse;
}

export interface ErrorResponse {
  error: string;
  code: string;
}

export interface HealthResponse {
  status: string;
  timestamp: string;
}

class ApiClient {
  private baseURL: string;

  constructor(baseURL: string = API_BASE_URL) {
    this.baseURL = baseURL;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`;

    const config: RequestInit = {
      headers: {
        "Content-Type": "application/json",
        ...options.headers,
      },
      ...options,
    };

    try {
      const response = await fetch(url, config);

      if (!response.ok) {
        try {
          const errorData: ErrorResponse = await response.json();
          throw new Error(
            errorData.error || `HTTP ${response.status}: ${response.statusText}`
          );
        } catch {
          throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }
      }

      // 204 No Contentの場合は空のオブジェクトを返す
      if (response.status === 204) {
        return {} as T;
      }

      return await response.json();
    } catch (error) {
      if (error instanceof Error) {
        throw error;
      }
      throw new Error("不明なエラーが発生しました");
    }
  }

  // 認証関連 - バックエンドのエンドポイントに合わせる
  async login(credentials: LoginRequest): Promise<AuthResponse> {
    return this.request<AuthResponse>("/api/v1/auth/login", {
      method: "POST",
      body: JSON.stringify(credentials),
    });
  }

  async register(userData: RegisterRequest): Promise<AuthResponse> {
    return this.request<AuthResponse>("/api/v1/auth/register", {
      method: "POST",
      body: JSON.stringify(userData),
    });
  }

  // ユーザー関連（認証が必要）
  async getUsers(
    page: number = 1,
    limit: number = 20,
    token?: string
  ): Promise<UsersResponse> {
    const headers: Record<string, string> = {};
    if (token) {
      headers.Authorization = `Bearer ${token}`;
    }

    return this.request<UsersResponse>(
      `/api/v1/users?page=${page}&limit=${limit}`,
      {
        method: "GET",
        headers,
      }
    );
  }

  async getUser(id: number, token?: string): Promise<User> {
    const headers: Record<string, string> = {};
    if (token) {
      headers.Authorization = `Bearer ${token}`;
    }

    return this.request<User>(`/api/v1/users/${id}`, {
      method: "GET",
      headers,
    });
  }

  async updateUser(
    id: number,
    userData: { email?: string; name?: string },
    token?: string
  ): Promise<User> {
    const headers: Record<string, string> = {};
    if (token) {
      headers.Authorization = `Bearer ${token}`;
    }

    return this.request<User>(`/api/v1/users/${id}`, {
      method: "PUT",
      body: JSON.stringify(userData),
      headers,
    });
  }

  async deleteUser(id: number, token?: string): Promise<void> {
    const headers: Record<string, string> = {};
    if (token) {
      headers.Authorization = `Bearer ${token}`;
    }

    return this.request<void>(`/api/v1/users/${id}`, {
      method: "DELETE",
      headers,
    });
  }

  // ヘルスチェック - 正しいエンドポイント
  async healthCheck(): Promise<HealthResponse> {
    return this.request<HealthResponse>("/health");
  }
}

export const apiClient = new ApiClient();

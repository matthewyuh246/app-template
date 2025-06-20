import React, { useEffect, useState } from "react";
import { useRouter } from "next/router";
import { User } from "../lib/api";

interface AuthState {
  token: string | null;
  user: User | null;
  isAuthenticated: boolean;
}

// ローカルストレージのキー
const TOKEN_KEY = "auth_token";
const USER_KEY = "auth_user";

class AuthManager {
  // トークンを保存
  static setToken(token: string): void {
    if (typeof window !== "undefined") {
      localStorage.setItem(TOKEN_KEY, token);
    }
  }

  // トークンを取得
  static getToken(): string | null {
    if (typeof window === "undefined") return null;
    return localStorage.getItem(TOKEN_KEY);
  }

  // ユーザー情報を保存
  static setUser(user: User): void {
    if (typeof window !== "undefined") {
      localStorage.setItem(USER_KEY, JSON.stringify(user));
    }
  }

  // ユーザー情報を取得
  static getUser(): User | null {
    if (typeof window === "undefined") return null;
    const userStr = localStorage.getItem(USER_KEY);
    return userStr ? JSON.parse(userStr) : null;
  }

  // 認証状態を取得
  static getAuthState(): AuthState {
    const token = this.getToken();
    const user = this.getUser();
    return {
      token,
      user,
      isAuthenticated: Boolean(token && user),
    };
  }

  // ログアウト（認証情報をクリア）
  static logout(): void {
    if (typeof window !== "undefined") {
      localStorage.removeItem(TOKEN_KEY);
      localStorage.removeItem(USER_KEY);
    }
  }

  // 認証情報を保存
  static setAuth(token: string, user: User): void {
    this.setToken(token);
    this.setUser(user);
  }
}

// 認証が必要なページを保護するコンポーネント
interface AuthWrapperProps {
  children: React.ReactNode;
  requireAuth?: boolean;
}

export function AuthWrapper({
  children,
  requireAuth = false,
}: AuthWrapperProps) {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(true);
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  useEffect(() => {
    const authState = AuthManager.getAuthState();
    setIsAuthenticated(authState.isAuthenticated);
    setIsLoading(false);

    if (requireAuth && !authState.isAuthenticated) {
      router.push("/login");
    }
  }, [router, requireAuth]);

  if (requireAuth) {
    if (isLoading) {
      return (
        <div className="min-h-screen flex items-center justify-center">
          <div className="text-center">
            <div className="animate-spin w-12 h-12 border-4 border-blue-500 border-t-transparent rounded-full mx-auto mb-4"></div>
            <p className="text-gray-600">認証を確認中...</p>
          </div>
        </div>
      );
    }

    if (!isAuthenticated) {
      return null; // リダイレクト中
    }
  }

  return <>{children}</>;
}

// 認証状態をチェックするカスタムフック
export function useAuth() {
  const [authState, setAuthState] = useState<AuthState>({
    token: null,
    user: null,
    isAuthenticated: false,
  });
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const state = AuthManager.getAuthState();
    setAuthState(state);
    setIsLoading(false);
  }, []);

  const login = (token: string, user: User) => {
    AuthManager.setAuth(token, user);
    setAuthState({ token, user, isAuthenticated: true });
  };

  const logout = () => {
    AuthManager.logout();
    setAuthState({ token: null, user: null, isAuthenticated: false });
  };

  return {
    ...authState,
    isLoading,
    login,
    logout,
  };
}

export { AuthManager };

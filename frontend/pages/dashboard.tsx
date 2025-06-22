import React, { useState, useEffect } from "react";
import { GetServerSideProps } from "next";
import Layout from "../components/Layout";
import { NextPageWithLayout } from "./_app";
import { AuthWrapper, useAuth } from "../components/AuthWrapper";
import { apiClient, User } from "../lib/api";

interface DashboardPageProps {
  initialUsers?: User[];
  error?: string;
}

const DashboardPage: NextPageWithLayout<DashboardPageProps> = ({
  initialUsers = [],
  error,
}) => {
  const { user, logout } = useAuth();
  const [users, setUsers] = useState<User[]>(initialUsers);
  const [loading, setLoading] = useState(false);
  const [searchTerm, setSearchTerm] = useState("");
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);

  // ユーザー一覧を取得
  const fetchUsers = async (page = 1, search = "") => {
    setLoading(true);
    try {
      const response = await apiClient.getUsers(page, 10, search);
      setUsers(response.users);
      setCurrentPage(response.pagination.page);
      setTotalPages(response.pagination.total_pages);
    } catch (error) {
      console.error("ユーザーの取得に失敗:", error);
    } finally {
      setLoading(false);
    }
  };

  // 検索処理
  const handleSearch = () => {
    fetchUsers(1, searchTerm);
  };

  // ログアウト処理
  const handleLogout = () => {
    logout();
    window.location.href = "/login";
  };

  return (
    <AuthWrapper requireAuth={true}>
      <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
        <div className="container mx-auto px-6 py-8">
          {/* ヘッダー */}
          <div className="glass-card rounded-2xl p-8 mb-8">
            <div className="flex flex-col md:flex-row justify-between items-start md:items-center space-y-4 md:space-y-0">
              <div>
                <h1 className="text-3xl font-bold gradient-text mb-2">
                  ダッシュボード
                </h1>
                <p className="text-gray-600">
                  こんにちは、{user?.name || "ゲスト"}さん！
                </p>
              </div>
              <button
                onClick={handleLogout}
                className="btn-secondary px-6 py-3 text-sm"
              >
                ログアウト
              </button>
            </div>
          </div>

          {/* エラー表示 */}
          {error && (
            <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded-lg mb-6">
              {error}
            </div>
          )}

          {/* 検索セクション */}
          <div className="glass-card rounded-2xl p-6 mb-8">
            <h2 className="text-xl font-bold mb-4">ユーザー検索</h2>
            <div className="flex space-x-4">
              <input
                type="text"
                placeholder="ユーザー名で検索..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="flex-1 px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                onKeyPress={(e) => e.key === "Enter" && handleSearch()}
              />
              <button
                onClick={handleSearch}
                disabled={loading}
                className="btn-primary px-6 py-3"
              >
                {loading ? "検索中..." : "検索"}
              </button>
            </div>
          </div>

          {/* ユーザー一覧 */}
          <div className="glass-card rounded-2xl p-6">
            <div className="flex justify-between items-center mb-6">
              <h2 className="text-xl font-bold">ユーザー一覧</h2>
              <p className="text-gray-600">{users.length}件のユーザー</p>
            </div>

            {loading ? (
              <div className="flex justify-center items-center py-12">
                <div className="animate-spin w-12 h-12 border-4 border-blue-500 border-t-transparent rounded-full"></div>
              </div>
            ) : users.length > 0 ? (
              <div className="space-y-4">
                {users.map((userData) => (
                  <div
                    key={userData.id}
                    className="bg-white/50 backdrop-blur-sm rounded-xl p-6 border border-white/20 hover:bg-white/70 transition-all duration-300"
                  >
                    <div className="flex flex-col md:flex-row justify-between items-start md:items-center space-y-2 md:space-y-0">
                      <div>
                        <h3 className="text-lg font-semibold text-gray-800">
                          {userData.name}
                        </h3>
                        <p className="text-gray-600">{userData.email}</p>
                      </div>
                      <div className="text-sm text-gray-500">
                        登録日:{" "}
                        {new Date(userData.created_at).toLocaleDateString(
                          "ja-JP"
                        )}
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="text-center py-12">
                <div className="text-gray-500 text-lg">
                  ユーザーが見つかりません
                </div>
                <p className="text-gray-400 mt-2">
                  検索条件を変更してみてください
                </p>
              </div>
            )}

            {/* ページネーション */}
            {totalPages > 1 && (
              <div className="flex justify-center items-center mt-8 space-x-2">
                <button
                  onClick={() => fetchUsers(currentPage - 1, searchTerm)}
                  disabled={currentPage === 1 || loading}
                  className="px-4 py-2 rounded-lg bg-white/70 text-gray-700 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-white transition-colors"
                >
                  前へ
                </button>
                <span className="px-4 py-2 text-gray-700">
                  {currentPage} / {totalPages}
                </span>
                <button
                  onClick={() => fetchUsers(currentPage + 1, searchTerm)}
                  disabled={currentPage === totalPages || loading}
                  className="px-4 py-2 rounded-lg bg-white/70 text-gray-700 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-white transition-colors"
                >
                  次
                </button>
              </div>
            )}
          </div>
        </div>
      </div>
    </AuthWrapper>
  );
};

// Pages RouterのgetServerSidePropsを使用してサーバーサイドで初期データを取得
export const getServerSideProps: GetServerSideProps<
  DashboardPageProps
> = async (context) => {
  try {
    // サーバーサイドでの認証チェック（トークンがCookieに保存されている場合）
    const token = context.req.cookies.auth_token;

    if (!token) {
      return {
        redirect: {
          destination: "/login",
          permanent: false,
        },
      };
    }

    // 初期データの取得を試行（オプション）
    // ここではクライアントサイドで取得することにする
    return {
      props: {
        initialUsers: [],
      },
    };
  } catch (error) {
    console.error("サーバーサイドエラー:", error);
    return {
      props: {
        initialUsers: [],
        error: "データの取得に失敗しました",
      },
    };
  }
};

// ページにレイアウトを適用（Pages Routerの機能）
DashboardPage.getLayout = function getLayout(page: React.ReactElement) {
  return <Layout>{page}</Layout>;
};

export default DashboardPage;

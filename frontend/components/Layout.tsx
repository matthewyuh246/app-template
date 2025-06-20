import React, { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/router";
import * as HeroIcons from "react-icons/hi";
import * as FontAwesome from "react-icons/fa";
import { AuthManager } from "./AuthWrapper";

interface LayoutProps {
  children: React.ReactNode;
}

export default function Layout({ children }: LayoutProps) {
  const {
    HiLightningBolt,
    HiHome,
    HiViewGrid,
    HiLogin,
    HiLogout,
    HiMenu,
    HiX,
    HiUserAdd,
  } = HeroIcons;
  const { FaGithub, FaTwitter } = FontAwesome;

  const router = useRouter();
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  // 認証状態チェック（クライアントサイドのみ）
  React.useEffect(() => {
    const authState = AuthManager.getAuthState();
    setIsAuthenticated(authState.isAuthenticated);
  }, []);

  const handleLogout = () => {
    AuthManager.logout();
    setIsAuthenticated(false);
    router.push("/");
  };

  const navigation = [
    { name: "ホーム", href: "/", icon: HiHome },
    {
      name: "ダッシュボード",
      href: "/dashboard",
      icon: HiViewGrid,
      requireAuth: true,
    },
  ];

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 via-white to-blue-50">
      {/* ナビゲーションヘッダー */}
      <header className="sticky top-0 z-50 glass-effect border-b border-white/20 backdrop-blur-md">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            {/* ロゴ */}
            <div className="flex-shrink-0">
              <Link href="/" className="group">
                <div className="flex items-center space-x-3">
                  <div className="w-10 h-10 bg-gradient-to-r from-blue-600 to-purple-600 rounded-xl flex items-center justify-center group-hover:scale-110 transition-transform duration-200 shadow-lg">
                    <HiLightningBolt className="w-5 h-5 text-white" />
                  </div>
                  <h1 className="text-xl font-bold gradient-text group-hover:scale-105 transition-transform duration-200">
                    App Template
                  </h1>
                </div>
              </Link>
            </div>

            {/* デスクトップナビゲーション */}
            <nav className="hidden md:flex space-x-1">
              {navigation.map((item) => {
                if (item.requireAuth && !isAuthenticated) return null;
                return (
                  <Link
                    key={item.name}
                    href={item.href}
                    className={`btn-ghost hover:bg-blue-50/50 hover:text-blue-600 ${
                      router.pathname === item.href
                        ? "bg-blue-50/50 text-blue-600"
                        : ""
                    }`}
                  >
                    <div className="flex items-center space-x-2">
                      <item.icon className="w-4 h-4" />
                      <span>{item.name}</span>
                    </div>
                  </Link>
                );
              })}
            </nav>

            {/* 認証ボタン（デスクトップ） */}
            <div className="hidden md:flex items-center space-x-3">
              {isAuthenticated ? (
                <button onClick={handleLogout} className="btn-secondary">
                  <div className="flex items-center space-x-2">
                    <HiLogout className="w-4 h-4" />
                    <span>ログアウト</span>
                  </div>
                </button>
              ) : (
                <div className="flex items-center space-x-3">
                  <Link href="/login" className="btn-ghost">
                    <div className="flex items-center space-x-2">
                      <HiLogin className="w-4 h-4" />
                      <span>ログイン</span>
                    </div>
                  </Link>
                  <Link href="/register" className="btn-primary">
                    <div className="flex items-center space-x-2">
                      <HiUserAdd className="w-4 h-4" />
                      <span>新規登録</span>
                    </div>
                  </Link>
                </div>
              )}
            </div>

            {/* モバイルメニューボタン */}
            <div className="md:hidden">
              <button
                onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
                className="p-2 rounded-lg hover:bg-gray-100/50 transition-colors duration-200"
              >
                {mobileMenuOpen ? (
                  <HiX className="w-5 h-5 text-gray-600" />
                ) : (
                  <HiMenu className="w-5 h-5 text-gray-600" />
                )}
              </button>
            </div>
          </div>
        </div>

        {/* モバイルメニュー */}
        {mobileMenuOpen && (
          <div className="md:hidden glass-effect border-t border-white/20">
            <div className="px-4 py-3 space-y-1">
              {navigation.map((item) => {
                if (item.requireAuth && !isAuthenticated) return null;
                return (
                  <Link
                    key={item.name}
                    href={item.href}
                    className={`block btn-ghost hover:bg-blue-50/50 hover:text-blue-600 ${
                      router.pathname === item.href
                        ? "bg-blue-50/50 text-blue-600"
                        : ""
                    }`}
                    onClick={() => setMobileMenuOpen(false)}
                  >
                    <div className="flex items-center space-x-2">
                      <item.icon className="w-4 h-4" />
                      <span>{item.name}</span>
                    </div>
                  </Link>
                );
              })}

              <div className="pt-2 mt-2 border-t border-white/20">
                {isAuthenticated ? (
                  <button
                    onClick={() => {
                      handleLogout();
                      setMobileMenuOpen(false);
                    }}
                    className="w-full text-left btn-ghost hover:bg-red-50/50 hover:text-red-600"
                  >
                    <div className="flex items-center space-x-2">
                      <HiLogout className="w-4 h-4" />
                      <span>ログアウト</span>
                    </div>
                  </button>
                ) : (
                  <div className="space-y-1">
                    <Link
                      href="/login"
                      className="block btn-ghost hover:bg-blue-50/50 hover:text-blue-600"
                      onClick={() => setMobileMenuOpen(false)}
                    >
                      <div className="flex items-center space-x-2">
                        <HiLogin className="w-4 h-4" />
                        <span>ログイン</span>
                      </div>
                    </Link>
                    <Link
                      href="/register"
                      className="block btn-primary"
                      onClick={() => setMobileMenuOpen(false)}
                    >
                      <div className="flex items-center space-x-2">
                        <HiUserAdd className="w-4 h-4" />
                        <span>新規登録</span>
                      </div>
                    </Link>
                  </div>
                )}
              </div>
            </div>
          </div>
        )}
      </header>

      {/* メインコンテンツ */}
      <main className="relative">
        {/* 背景装飾 */}
        <div className="fixed inset-0 overflow-hidden pointer-events-none">
          <div className="absolute -top-40 -right-40 w-80 h-80 bg-purple-300 rounded-full mix-blend-multiply filter blur-xl opacity-20 floating-animation"></div>
          <div
            className="absolute -bottom-40 -left-40 w-80 h-80 bg-blue-300 rounded-full mix-blend-multiply filter blur-xl opacity-20 floating-animation"
            style={{ animationDelay: "2s" }}
          ></div>
          <div
            className="absolute top-40 left-1/2 transform -translate-x-1/2 w-80 h-80 bg-pink-300 rounded-full mix-blend-multiply filter blur-xl opacity-20 floating-animation"
            style={{ animationDelay: "4s" }}
          ></div>
        </div>

        {children}
      </main>

      {/* フッター */}
      <footer className="relative mt-20 glass-effect border-t border-white/20">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
          <div className="text-center">
            <div className="flex justify-center items-center space-x-3 mb-6">
              <div className="w-8 h-8 bg-gradient-to-r from-blue-600 to-purple-600 rounded-xl flex items-center justify-center shadow-lg">
                <HiLightningBolt className="w-4 h-4 text-white" />
              </div>
              <span className="text-xl font-bold gradient-text">
                App Template
              </span>
            </div>
            <p className="text-gray-600 text-sm mb-6 max-w-md mx-auto">
              Next.js、Go、TailwindCSSで構築されたモダンなWebアプリケーションテンプレート
            </p>
            <div className="flex justify-center space-x-6">
              <a
                href="https://github.com"
                target="_blank"
                rel="noopener noreferrer"
                className="text-gray-400 hover:text-blue-600 transition-colors duration-200"
              >
                <span className="sr-only">GitHub</span>
                <FaGithub className="w-6 h-6" />
              </a>
              <a
                href="https://twitter.com"
                target="_blank"
                rel="noopener noreferrer"
                className="text-gray-400 hover:text-blue-600 transition-colors duration-200"
              >
                <span className="sr-only">Twitter</span>
                <FaTwitter className="w-6 h-6" />
              </a>
            </div>
            <div className="mt-8 pt-6 border-t border-white/20">
              <p className="text-gray-500 text-xs">
                © 2024 App Template. All rights reserved.
              </p>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}

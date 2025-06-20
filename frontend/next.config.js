/** @type {import('next').NextConfig} */
const nextConfig = {
  // 本番環境では静的エクスポートを使用
  output: process.env.NODE_ENV === "production" ? "export" : undefined,
  trailingSlash: false,
  images: {
    unoptimized: true,
  },
  async rewrites() {
    // 開発環境でのみrewritesを使用（静的エクスポート時は無効）
    if (process.env.NODE_ENV === "development") {
      const apiBaseUrl =
        process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080";
      return [
        {
          source: "/api/:path*",
          destination: `${apiBaseUrl}/:path*`,
        },
      ];
    }
    return [];
  },
};

module.exports = nextConfig;

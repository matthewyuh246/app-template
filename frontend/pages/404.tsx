import React from "react";
import Link from "next/link";
import Layout from "../components/Layout";
import { NextPageWithLayout } from "./_app";

const Custom404: NextPageWithLayout = () => {
  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center px-4 sm:px-6 lg:px-8">
      <div className="max-w-lg w-full text-center">
        <div className="glass-card rounded-2xl p-12">
          {/* 404イラスト */}
          <div className="mb-8">
            <div className="text-8xl font-bold gradient-text mb-4">404</div>
            <div className="w-24 h-1 bg-gradient-to-r from-blue-500 to-purple-600 mx-auto rounded-full"></div>
          </div>

          {/* メッセージ */}
          <h1 className="text-2xl md:text-3xl font-bold text-gray-800 mb-4">
            ページが見つかりません
          </h1>
          <p className="text-gray-600 mb-8 leading-relaxed">
            お探しのページは存在しないか、移動した可能性があります。
            <br />
            URLをご確認いただくか、ホームページからご利用ください。
          </p>

          {/* アクションボタン */}
          <div className="space-y-4">
            <Link
              href="/"
              className="inline-block w-full btn-primary py-3 px-6 text-lg rounded-lg"
            >
              ホームページに戻る
            </Link>
            <button
              onClick={() => window.history.back()}
              className="inline-block w-full btn-secondary py-3 px-6 text-lg rounded-lg"
            >
              前のページに戻る
            </button>
          </div>
        </div>

        {/* 装飾要素 */}
        <div className="absolute top-1/4 left-1/4 w-32 h-32 bg-blue-200 rounded-full opacity-20 floating"></div>
        <div
          className="absolute bottom-1/4 right-1/4 w-24 h-24 bg-purple-200 rounded-full opacity-20 floating"
          style={{ animationDelay: "2s" }}
        ></div>
      </div>
    </div>
  );
};

// 独立したレイアウト（エラーページ用）
Custom404.getLayout = function getLayout(page: React.ReactElement) {
  return page; // レイアウトなしで表示
};

export default Custom404;

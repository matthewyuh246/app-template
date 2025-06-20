import React from "react";
import type { AppProps } from "next/app";
import type { NextPage } from "next";
import Head from "next/head";
import { Inter } from "next/font/google";
import "styles/globals.css";

const inter = Inter({ subsets: ["latin"] });

// ページコンポーネントにレイアウト情報を追加するための型定義
export type NextPageWithLayout<P = {}, IP = P> = NextPage<P, IP> & {
  getLayout?: (page: React.ReactElement) => React.ReactNode;
  requireAuth?: boolean;
};

type AppPropsWithLayout = AppProps & {
  Component: NextPageWithLayout;
};

export default function MyApp({ Component, pageProps }: AppPropsWithLayout) {
  // ページごとのレイアウト設定を取得（なければデフォルトでそのまま返す）
  const getLayout = Component.getLayout ?? ((page) => page);

  return (
    <>
      <Head>
        <title>App Template</title>
        <meta
          name="description"
          content="モダンな技術スタックを使用したWebアプリケーションテンプレート"
        />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <div className={inter.className}>
        {getLayout(<Component {...pageProps} />)}
      </div>
    </>
  );
}

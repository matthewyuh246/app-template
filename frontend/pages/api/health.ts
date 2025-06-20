import type { NextApiRequest, NextApiResponse } from "next";

interface HealthResponse {
  status: string;
  timestamp: string;
  environment: string;
  version: string;
}

export default function handler(
  req: NextApiRequest,
  res: NextApiResponse<HealthResponse>
) {
  // GET メソッドのみ許可
  if (req.method !== "GET") {
    res.setHeader("Allow", ["GET"]);
    res.status(405).end(`Method ${req.method} Not Allowed`);
    return;
  }

  // ヘルスチェック情報を返す
  const healthData: HealthResponse = {
    status: "ok",
    timestamp: new Date().toISOString(),
    environment: process.env.NODE_ENV || "development",
    version: process.env.npm_package_version || "1.0.0",
  };

  res.status(200).json(healthData);
}

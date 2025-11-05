export type AppTokenCostsResponse = {
  data: Array<{
    date: string;
    token_count: number;
    total_price: number;
    currency: number;
  }>;
};

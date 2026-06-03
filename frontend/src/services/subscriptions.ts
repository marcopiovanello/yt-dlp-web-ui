// import { PaginatedResponse } from '../types'

export type Subscription = {
  id: string
  url: string
  params: string
  cron_expression: string
}

// class SubscriptionService {
//   private _baseURL: string = ''

//   public set baseURL(v: string) {
//     this._baseURL = v
//   }

//   public async delete(id: string): Promise<void> {

//   }

//   public async listPaginated(start: number, limit: number = 50): Promise<PaginatedResponse<Subscription[]>> {
//     const res = await fetch(`${this._baseURL}/subscriptions?id=${start}&limit=${limit}`)
//     const data: PaginatedResponse<Subscription[]> = await res.json()

//     return data
//   }

//   public async submit(sub: Subscription): Promise<void> {

//   }

//   public async edit(sub: Subscription): Promise<void> {

//   }
// }

// export default SubscriptionService
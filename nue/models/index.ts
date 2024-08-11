// file: models/index.ts
export default {

    async getFeedback(): Promise<any> {
      const resp = await fetch('/models/feedback.json')
      return await resp.json()
    },
  
    async getUsers(): Promise<any> {
      const resp = await fetch('/models/users.json')
      return await resp.json()
    },
  
    async getAnalytics(): Promise<any> {
      const resp = await fetch('/models/analytics.json')
      return await resp.json()
    }
  }
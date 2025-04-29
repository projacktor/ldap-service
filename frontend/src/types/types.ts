export interface User {
  preferred_username: string
  email: string
  resource_access: {
    account: {
      roles: string[]
    }
  }
}

export interface User {
  id: string;
  user_name: string;
  first_name: string;
  last_name: string;
  email: string;
  department?: string;
}

export interface UserSnowflake {
  id: number;
  userName: string;
  firstName: string;
  lastName: string;
  email: string;
  department?: string;
}
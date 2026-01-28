export interface User {
  id: string | number;
  user_name: string;
  first_name: string;
  last_name: string;
  email: string;
  department?: string;
}
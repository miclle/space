import axios, { AxiosRequestConfig } from 'axios';

export class HTTP {
  GET = async <T = any>(url: string, config?: AxiosRequestConfig): Promise<T> => {
    const resp = await axios.get<T>(url, config);
    return resp.data;
  };

  HEAD = axios.head;

  OPTIONS = axios.options;

  POST = async <T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> => {
    const resp = await axios.post<T>(url, data, config);
    return resp.data;
  };

  PUT = async <T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> => {
    const resp = await axios.put<T>(url, data, config);
    return resp.data;
  };

  PATCH = async <T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> => {
    const resp = await axios.patch<T>(url, data, config);
    return resp.data;
  };

  DELETE = async <T = any>(url: string, config?: AxiosRequestConfig): Promise<T> => {
    const resp = await axios.delete<T>(url, config);
    return resp.data;
  };
}

const instance = new HTTP();

export default instance;
export const { GET, POST, PUT, PATCH, DELETE } = instance;

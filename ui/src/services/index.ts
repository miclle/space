import axios, { AxiosError } from 'axios';
import { map } from 'lodash';
import { notification } from 'antd';

export type { AxiosError, AxiosResponse } from 'axios';

export * from 'services/lib/http';

export interface IErrorMessage<T = any> {
  code: string
  message: T
}

axios.defaults.baseURL = '/api';
axios.defaults.withCredentials = true;

// Add a response interceptor
/* eslint arrow-body-style: "off" */
axios.interceptors.response.use((response) => {
  return response;
}, (error: AxiosError<IErrorMessage>) => {

  if (error.response?.status === 403) {
    notification.error({
      key: 'ErrForbidden',
      message: 'Forbidden',
      description: map(error.response?.data.message, (value, key) => value).join('\n')
    });
  }

  return Promise.reject(error.response);
});

// pagination
export * from './pagination';

export * as Account from './account';
export * as Space from './space';
export * as Page from './page';
export * as Markdown from './markdown';

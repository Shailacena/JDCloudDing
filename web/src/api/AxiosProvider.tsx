import { createContext, useContext, useState } from "react";
import { request } from "./request";
import { AxiosError, AxiosInstance, AxiosResponse } from "axios";
import { AUTH_TYPE, useAppContext } from "../AppProvider";
import { message } from "antd";

interface AxiosInstanceContextType {
  axiosInstance: AxiosInstance;
  get: (url: string, params?: any) => Promise<AxiosResponse | void>
  post: (url: string, data: any) => Promise<AxiosResponse | void>
}

const AxiosContext = createContext<AxiosInstanceContextType>(null!)

export function AxiosProvider({ children }: { children: React.ReactNode }) {
  const [axiosInstance] = useState(() => request);
  let app = useAppContext()
  let headers: any = {}
  let path = app.cookie.path

  if (app?.cookie?.token) {
    headers["Ttttt"] = app.cookie.token
  }

  if (app?.cookie?.id) {
    headers["Ddd"] = app.cookie.id
  }

  if (app?.cookie?.role) {
    headers["Rrrr"] = app.cookie.role
  }

  let get = (url: string, params?: any): Promise<AxiosResponse | void> => {
    return request.get(url, { headers, params: params }).then((res) => {
      return res
    }).catch((e) => {
      if (e.status == 401) {
        handle401(path)
      }
    })
  }

  let post = (url: string, data: any): Promise<AxiosResponse | void> => {
    return request.post(url, data, { headers }).then((res) => {
      return res
    }).catch((e: AxiosError) => {
      if (e.status == 401) {
        handle401(path)
      } else if (e.status == 500) {
        //@ts-ignore
        if (e.response && e.response.data && e.response.data.message)
          //@ts-ignore
          message.error(e.response.data.message)
      } else if (e.status == 403) {
        //@ts-ignore
        if (e.response && e.response.data && e.response.data.message)
          //@ts-ignore
          message.error(e.response.data.message)
      }
      throw e;
    })
  }

  function handle401(path?: AUTH_TYPE) {
    switch (path) {
      case AUTH_TYPE.ADMIN:
        app.auth?.adminSignout(() => { })
        break;

        case AUTH_TYPE.PARTNER:
          app.auth?.partnerSignout(() => { })
        break;

        case AUTH_TYPE.MERCHANT:
          app.auth?.merchantSignout(() => { })
        break;

      default:
        console.error("handle 401 error");
        break;
    }
  }

  let value = {
    axiosInstance,
    get,
    post,
  }

  return <AxiosContext.Provider value={value} >
    {children}
  </AxiosContext.Provider>
}

export function useAxios() {
  // 获取axios实例
  return useContext(AxiosContext)
}


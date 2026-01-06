import React from "react";
import { Navigate, useLocation } from "react-router-dom";
import { useCookies } from "react-cookie";
import { AdminLoginResp, IMerchant, IPartner, MerchantLoginResp, PartnerLoginResp } from "./api/types";
import { getExpirationDate } from "./utils/Tool";

export enum AUTH_TYPE {
  ADMIN = '/web/admin',
  PARTNER = '/web/partner',
  MERCHANT = '/web/merchant'
}

enum CookieFiled {
  Token = 'token',
  Nickname = 'nickname',
  Username = 'username',
  Role = 'role',
  ID = 'id',
  LEVEL = 'level',
  PATH = 'path'
}

interface Cookie {
  token?: string
  nickname?: string
  username?: string;
  role?: number
  id?: number
  level?: number;
  path?: AUTH_TYPE;
}

interface AuthContextType {
  auth: {
    adminSignin: (value: AdminLoginResp, id: string, callback: Function) => void;
    adminSignout: (callback: Function) => void;

    partnerSignin: (value: PartnerLoginResp, id: string, callback: Function) => void;
    partnerSignout: (callback: Function) => void;

    merchantSignin: (value: MerchantLoginResp, id: string, callback: Function) => void;
    merchantSignout: (callback: Function) => void;
  }
  cookie: Cookie;

  partnerList: IPartner[];
  merchantList: IMerchant[];
}

let AuthContext = React.createContext<AuthContextType>(null!);

export function useAppContext() {
  return React.useContext(AuthContext);
}

function AppProvider({ children }: { children: React.ReactNode }) {
  let [partnerList] = React.useState<any>(null);
  let [merchantList] = React.useState<any>(null);

  let [cookie, setCookie, removeCookie] = useCookies([
    CookieFiled.Token,
    CookieFiled.Nickname,
    CookieFiled.Username,
    CookieFiled.Role,
    CookieFiled.ID,
    CookieFiled.LEVEL,
    CookieFiled.PATH
  ]);

  const adminPath = AUTH_TYPE.ADMIN
  const partnerPath = AUTH_TYPE.PARTNER
  const merchantPath = AUTH_TYPE.MERCHANT

  let adminSignin = async (data: AdminLoginResp, username: string, callback: Function) => {
    let exp = getExpirationDate(7)
    setCookie(CookieFiled.Token, data.token, { path: adminPath, expires: exp });
    setCookie(CookieFiled.Nickname, data.nickname, { path: adminPath, expires: exp });
    setCookie(CookieFiled.Role, data.role, { path: adminPath, expires: exp });
    setCookie(CookieFiled.Username, username, { path: adminPath, expires: exp });
    setCookie(CookieFiled.ID, data.id, { path: adminPath, expires: exp });
    setCookie(CookieFiled.LEVEL, 0, { path: adminPath, expires: exp });
    setCookie(CookieFiled.PATH, adminPath, { path: adminPath, expires: exp });
    callback()
  };

  let adminSignout = (callback: Function) => {
    removeCookie(CookieFiled.Token, { path: adminPath })
    removeCookie(CookieFiled.Nickname, { path: adminPath })
    removeCookie(CookieFiled.Role, { path: adminPath })
    removeCookie(CookieFiled.Username, { path: adminPath })
    removeCookie(CookieFiled.ID, { path: adminPath })
    removeCookie(CookieFiled.LEVEL, { path: adminPath })
    removeCookie(CookieFiled.PATH, { path: adminPath })

    callback();
  };

  let partnerSignin = async (data: PartnerLoginResp, username: string, callback: Function) => {
    let exp = getExpirationDate(7)
    setCookie(CookieFiled.Token, data.token, { path: partnerPath, expires: exp });
    setCookie(CookieFiled.Nickname, data.nickname, { path: partnerPath, expires: exp });
    setCookie(CookieFiled.Username, username, { path: partnerPath, expires: exp });
    setCookie(CookieFiled.ID, data.id, { path: partnerPath, expires: exp });
    setCookie(CookieFiled.LEVEL, 0, { path: partnerPath, expires: exp });
    setCookie(CookieFiled.PATH, partnerPath, { path: partnerPath, expires: exp });
    callback()
  };

  let partnerSignout = (callback: Function) => {
    removeCookie(CookieFiled.Token, { path: partnerPath })
    removeCookie(CookieFiled.Nickname, { path: partnerPath })
    removeCookie(CookieFiled.ID, { path: partnerPath })
    removeCookie(CookieFiled.Username, { path: adminPath })
    removeCookie(CookieFiled.LEVEL, { path: partnerPath })
    removeCookie(CookieFiled.PATH, { path: partnerPath })
    callback();
  };

  let merchantSignin = async (data: MerchantLoginResp, username: string, callback: Function) => {
    let exp = getExpirationDate(7)
    setCookie(CookieFiled.Token, data.token, { path: merchantPath, expires: exp });
    setCookie(CookieFiled.Nickname, data.nickname, { path: merchantPath, expires: exp });
    setCookie(CookieFiled.Username, username, { path: merchantPath, expires: exp });
    setCookie(CookieFiled.ID, data.id, { path: merchantPath, expires: exp });
    setCookie(CookieFiled.LEVEL, 0, { path: merchantPath, expires: exp });
    setCookie(CookieFiled.PATH, merchantPath, { path: merchantPath, expires: exp });
    callback()
  };

  let merchantSignout = (callback: Function) => {
    removeCookie(CookieFiled.Token, { path: merchantPath })
    removeCookie(CookieFiled.Nickname, { path: merchantPath })
    removeCookie(CookieFiled.ID, { path: merchantPath })
    removeCookie(CookieFiled.Username, { path: adminPath })
    removeCookie(CookieFiled.LEVEL, { path: merchantPath })
    removeCookie(CookieFiled.PATH, { path: merchantPath })
    callback();
  };

  let value = {
    auth: { adminSignin, adminSignout, partnerSignin, partnerSignout, merchantSignin, merchantSignout },
    cookie,
    partnerList,
    merchantList,
  }
  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function RequireAuth({ children }: { children: JSX.Element }) {
  // let ctx = useAppContext();

  let location = useLocation();
  let [cookies] = useCookies([CookieFiled.Token]);

  // 已登陆
  if (cookies.token) {
    // ctx.auth.token = cookies.token
    return children;
  }

  // 重定向至login页面，但是保存用户试图访问的location，这样我们可以把登陆后的用户重定向至那个页面
  return <Navigate to="/admin/login" state={{ from: location }} replace />;
}

export function RequireAuthPartner({ children }: { children: JSX.Element }) {
  // let ctx = useAppContext();
  let location = useLocation();
  let [cookies] = useCookies([CookieFiled.Token]);

  // 已登陆
  if (cookies.token) {
    // ctx.auth.token = cookies.token
    return children;
  }

  // 重定向至login页面，但是保存用户试图访问的location，这样我们可以把登陆后的用户重定向至那个页面
  return <Navigate to="/partner/login" state={{ from: location }} replace />;
}

export function RequireAuthMerchant({ children }: { children: JSX.Element }) {
  // let ctx = useAppContext();
  let location = useLocation();
  let [cookies] = useCookies([CookieFiled.Token]);

  // 已登陆
  if (cookies.token) {
    return children;
  }

  // 重定向至login页面，但是保存用户试图访问的location，这样我们可以把登陆后的用户重定向至那个页面
  return <Navigate to="/merchant/login" state={{ from: location }} replace />;
}

export default AppProvider
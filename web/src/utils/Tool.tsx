
/**
 * 获取随机数
 * @param start: number 随机区间头
 * @param end: number 随机区间尾
 * @returns value: number
 */
export const getRandomNumber = (start: number, end: number) => {
  return Math.floor(Math.random() * end) + start;
}

/**
 * 获取格式化日期
 * @param date: Date 日期
 * @returns value: string
 */
export const getDataFormat = (date: Date) => {
  return `${date.getFullYear()}-${('0' + (date.getMonth() + 1)).slice(-2)}-${('0' + date.getDate()).slice(-2)} ${('0' + date.getHours()).slice(-2)}:${('0' + date.getMinutes()).slice(-2)}:${('0' + date.getSeconds()).slice(-2)}`;
}

/**
 * 获取Cookies
 * @param name: string key
 * @returns value: string
 */
export const getCookieByDocument = (name: string) => {
  const cookies = document.cookie
    .split("; ")
    .find((row) => row.startsWith(`${name}=`));

  return cookies ? cookies.split("=")[1] : null;
};

/**
 * 设置Cookies
 * @param name: string key
 * @param value: string value
 * @param days: number 有效天数
 * @param path: string 有效path
 */
export const setCookieByDocument = (name: string, value: string, days = 1, path = '/') => {
  const expirationDate = new Date();
  expirationDate.setDate(expirationDate.getDate() + days);
  document.cookie = `${name}=${value}; expires=${expirationDate.toUTCString()}; path=${path}`;
};

/**
 * 获取有效期
 * @param days: number 有效天数
 */
export const getExpirationDate = (days = 1) => {
  const expirationDate = new Date();
  expirationDate.setDate(expirationDate.getDate() + days);
  return expirationDate;
};

/**
 * 获取cookies路径
 * @param name
 * @returns
 */
export const getCookiePath = (name: string) => {
  const cookiePattern = new RegExp(`(^|;\\s*)${name}=([^;]*)`);
  const cookieMatch = document.cookie.match(cookiePattern);
  if (cookieMatch) {
    const [pathString] = cookieMatch[0].split(';');
    const pathPattern = /path=([^;]*)/;
    const pathMatch = pathPattern.exec(pathString);
    return pathMatch ? pathMatch[1] : null;
  }
  return null;
};

/**
 * 获取随机路由
 * @returns path: string
 */
export const getRandomPath = (index: number) => {
  // 检查本地中是否存在路由缓存
  let paths: string[] = [];
  let pathstr = localStorage.getItem("paths");
  if (pathstr) {
    paths = paths.concat(pathstr.split(","));
    if (paths[index] && paths[index] != '') {
      return paths[index]
    }
  }

  const characters = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
  const charactersLength = characters.length;
  // 一次生成10个
  for (let i = 0; i < 10; i++) {
    // 随机路由4-8个字符
    let length = getRandomNumber(4, 8);
    let result = '';
    if (window.crypto && window.crypto.getRandomValues) {
      const randomValues = new Uint32Array(length);
      window.crypto.getRandomValues(randomValues);
      for (let i = 0; i < length; i++) {
        result += characters[randomValues[i] % charactersLength];
      }
    } else {
      // 使用 Math.random() 函数作为备选方法
      for (let i = 0; i < length; i++) {
        const randomIndex = Math.floor(Math.random() * charactersLength);
        result += characters[randomIndex];
      }
    }
    paths.push(result)
  }
  localStorage.setItem('paths', paths.toString());
  if (paths[index] && paths[index] != '') {
    return paths[index]
  } else {
    getRandomPath(index);
  }
}

/**
 * 根据时间问候
 * @returns
 */
export function greet(): string {
  const hours = new Date().getHours();
  if (!isNaN(hours)) {
    if (hours >= 6 && hours < 11) {
      return '早上好';
    } else if (hours >= 11 && hours < 13) {
      return '中午好';
    } else if (hours >= 13 && hours < 19) {
      return '下午好';
    } else if (hours >= 18 && hours < 24) {
      return '晚上好';
    } else if (hours >= 0 && hours < 6) {
      return '夜深了';
    }
    return '您好';
  } else {
    return '您好';
  }
}


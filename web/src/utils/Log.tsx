

export function log(...args: any[]){
    // 获取堆栈跟踪
    const stack = new Error().stack;
    if (stack) {
        const line2 = stack.split('\n')[1];
        const line2s = line2.split('/');
        console.log(line2s.pop(), args);
    }
}
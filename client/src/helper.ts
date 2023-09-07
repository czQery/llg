export const formatTime = (value: number): string => {
    let h = Math.trunc(value / 60);
    let m = value % 60;

    let hs: string = h.toString();
    let ms: string = m.toString();

    if (h.toString().length === 1) {
        hs = "0" + h.toString();
    }

    if (m.toString().length === 1) {
        ms = "0" + m.toString();
    }

    return hs + ":" + ms;
}
export const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms));
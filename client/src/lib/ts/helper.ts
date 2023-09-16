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

export const formatDate = (value: number): string => {
    let date = new Date(value * 60 * 60 * 24 * 1000);
    return date.getDate().toString()+"."+(date.getMonth()+1).toString();
}

export const formatDuration = (minutes: number): string => {
    if (minutes >= 60) {
        return (minutes / 60).toFixed(1).replace(".0", "") + "h";
    } else {
        return minutes.toFixed(1).replace(".0", "") + "m";
    }
}

export const getDate = (date: Date): string => {
    return date.getFullYear().toString() + "-" + ((date.getMonth() + 1).toString().length == 1 ? "0" + (date.getMonth() + 1).toString() : (date.getMonth() + 1).toString())
}

export const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms));

export const getRandomColor = (): string => {
    let c = Math.floor(Math.random() * 16777215).toString(16);

    if (c.length != 6) {
        c = c+"0";
    }

    return "#" + c;
}

export const getPaletteColor = (): string => {
    const colors: string[] = ["f28682", "0ceaf9", "fff654", "a845b3", "a09262", "d5334b", "60dd83", "f5980c", "fa147f"]
    return "#"+colors[Math.floor(Math.random() * (colors.length - 1))]
}

export const isJSON = (response: Response): boolean => {
    return response.headers.get("content-type")!.includes("application/json");
}

export const isOK = (data: any): boolean => {
    return !(data == null);
}
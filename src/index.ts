import { app, BrowserWindow } from "electron";
import installExtension, { VUEJS_DEVTOOLS } from "electron-devtools-installer";
import path from "path";

let mainWindow: BrowserWindow | null;

const createMainWindow = () => {
    mainWindow = new BrowserWindow({
        backgroundColor: "#0d1122",
        height: 600,
        width: 960,
    });

    mainWindow.loadFile(path.join(__dirname, "index.html"));

    if (process.env.NODE_ENV !== "production") {
        installExtension(VUEJS_DEVTOOLS);
    }

    mainWindow.on("closed", () => mainWindow = null);
};

app.on("ready", createMainWindow);

app.on("window-all-closed", () => {
    if (process.platform !== "darwin") {
        app.quit();
    }
});

app.on("activate", () => {
    if (mainWindow === null) {
        createMainWindow();
    }
});

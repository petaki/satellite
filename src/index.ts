import { app, BrowserWindow } from "electron";
import path from "path";

let mainWindow: BrowserWindow | null;

const createMainWindow = () => {
    mainWindow = new BrowserWindow({
        height: 600,
        width: 800,
    });

    mainWindow.loadFile(path.join(__dirname, "index.html"));
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

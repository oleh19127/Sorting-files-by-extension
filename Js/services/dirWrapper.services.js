import { existsSync, promises } from "fs";
import { print } from "./print.services.js";

class DirWrapper {
  async createDirIfNotExist(path) {
    if (!existsSync(path)) {
      await promises.mkdir(path, { recursive: true });
      return print.warning(`Create: ${path}`);
    }
  }

  async deleteDir(dir) {
    await promises.rmdir(dir);
    return print.error(`Folder removed: ${dir}`);
  }
}

const dirWrapper = new DirWrapper();

export { dirWrapper };

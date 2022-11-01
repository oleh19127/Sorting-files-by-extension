import { print } from "./print.services.js";
import { dirname, join, extname, basename } from "path";
import { existsSync, promises } from "fs";

class FileWrapper {
  async moveFile(oldPath, newPath) {
    const checkPath = this.incrementFile(newPath);
    await promises.rename(oldPath, checkPath);
    print.successful(`${oldPath} >> ${checkPath}`);
  }

  incrementFile(filepath) {
    while (existsSync(filepath)) {
      let dir = dirname(filepath);
      let ext = extname(filepath);
      let base = basename(filepath, ext);

      let re = /\((\d+)\)$/;
      let i = re.exec(base);

      if (i && i[1]) {
        i = i[1];
        i = parseInt(i) + 1;
        base = base.replace(re, "(" + i + ")");
      } else {
        base += "(1)";
      }
      filepath = join(dir, base + ext);
      if (!existsSync(filepath)) {
        return filepath;
      }
    }
    return filepath;
  }
}

const fileWrapper = new FileWrapper();

export { fileWrapper };

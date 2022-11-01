import { structure } from "../storage/structure.storage.js";
import { basename, extname, join } from "path";
import { globby } from "globby";
import { statSync, promises } from "fs";
import { dirWrapper } from "./dirWrapper.services.js";
import { fileWrapper } from "./fileWrapper.services.js";

class Sort {
  constructor(structure) {
    this.structure = structure;
  }

  async byExtension() {
    const arrByExt = await this.checkFiles();
    for (const element of arrByExt) {
      await dirWrapper.createDirIfNotExist(element.folder);
      await fileWrapper.moveFile(element.oldPath, element.newPath);
    }
  }

  async removeAllEmptyFolders() {
    const allFiles = await globby(
      ["**/**/*", "!services", "!storage", "!node_modules", "!Sorted files"],
      { onlyDirectories: true, dot: true }
    );
    for (const dir of allFiles.reverse()) {
      const filesInDir = await promises.readdir(dir);
      if (filesInDir.length <= 0) {
        await dirWrapper.deleteDir(dir);
      }
    }
  }

  async checkFiles() {
    const allFiles = await globby(
      ["**/**/*", "!services", "!storage", "!node_modules", "!Sorted files"],
      { onlyFiles: true, dot: true }
    );
    let paths = [];
    allFiles.forEach((path) => {
      for (const key in this.structure) {
        if (this.structure.hasOwnProperty(key)) {
          const data = this.structure[key];
          data.extensions.forEach((extension) => {
            const file = basename(path);
            const file_ext = extname(file).replace(".", "");
            if (extension.toLowerCase() === file_ext.toLowerCase()) {
              const modTimeYear = statSync(path).mtime.getFullYear().toString();
              const folder = join("Sorted files", data.folder, modTimeYear);
              const newPath = join(folder, file);
              paths.push({ oldPath: path, newPath, folder });
            }
          });
        }
      }
    });
    return paths;
  }
}

const sort = new Sort(structure);

export { sort };

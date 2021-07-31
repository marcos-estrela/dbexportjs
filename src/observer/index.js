// require('../config')

const fs = require("fs");
const sourceFile = "../database/logs/general-log.log";
const destinationFile = "../database/logs/general-log-copy.log";

class DbObserver {
  constructor() {
    this.lastCalledTime = new Date();
  }

  async changeTriggered(date) {
    try {
      const timeDiff = Math.abs(date.getTime() - this.lastCalledTime.getTime());
      this.lastCalledTime = date;
      if (timeDiff > 2000 && !this.checkFileIsEmpty()) {
        this.handleGeneralLog();
      }
    } catch (err) {
      console.error(err);
    }
  }

  async handleGeneralLog() {
    try {
      while (
        Math.abs(new Date().getTime() - this.lastCalledTime.getTime()) < 2000
      ) {}

      console.info("Handling general log...");

      // this.createGeneralLogCopy();
      // this.cleanGeneralLog();
      this.processGeneralLog();
    } catch (err) {
      console.error(err);
    }
  }

  /**
   * Creates a copy from the original general_log file
   *
   */
  createGeneralLogCopy() {
    try {
      process.stdout.write("Creating a copy from general_log... ");
      fs.copyFileSync(sourceFile, destinationFile);
      process.stdout.write("[OK]\n");
    } catch (err) {
      console.error(err);
    }
  }

  /**
   * Clean the original general_log file
   *
   */
  cleanGeneralLog() {
    try {
      process.stdout.write("Cleaning general_log... ");
      fs.writeFileSync(sourceFile, "");
      process.stdout.write("[OK]\n");
    } catch (err) {
      console.error(err);
    }
  }

  processGeneralLog() {
    try {
      console.log("Processing general log...");
      const data = fs.readFileSync(destinationFile, "UTF-8");
      const lines = data.split(/\r?\n/);
      lines.forEach((line) => {
        this.processGeneralLogLine(line);
      });
    } catch (err) {
      console.error(err);
    }
  }

  processGeneralLogLine(line) {
    console.info(line);
    // TODO: Identifica os padrões nas linhas
  }

  checkFileIsEmpty() {
    try {
      const content = fs.readFileSync(sourceFile);
      if (content.length) {
        console.warn("File has content.");
        return false;
      } else {
        console.warn("File is empty.");
        return true;
      }
    } catch (err) {
      console.error(err);
    }
  }
}

module.exports = DbObserver;

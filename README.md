# ✨ xrun ✨

A utility tool to streamline running multiple project commands effortlessly in a single terminal session! Perfect for managing development environments with ease.

---


## ⚡ Features

- Define and manage multiple commands in a single JSON configuration file.
- Execute commands simultaneously with ease.
- Works with any project—Laravel, Vue, Node.js, etc.

---

## 🔧 Installation (Linux)

### 🔄 Download

```bash
curl -s https://api.github.com/repos/sanda0/xrun/releases/latest \
| grep "browser_download_url" \
| grep "xrun-linux-amd64" \
| cut -d '"' -f 4 \
| xargs curl -L -o xrun
```

### 🔧 Install

```bash
sudo chmod +x xrun && sudo mv xrun /usr/bin/xrun
```

---

## 🚀 How to Use

### 🔍 Step 1: Go to Your Project Folder

```bash
cd /<path_to_your_project>/project
```

### 🔄 Step 2: Initialize Configuration

Run the init command to generate the configuration file:

```bash
xrun --init
```

You will see a file named `config.xrun.json` in your project directory. Open it to edit the configuration.

### 🗄 Step 3: Edit Configuration File

Here’s an example of how your `config.xrun.json` might look:

```json
{
  "Commands": [
    {
      "Label": "Laravel",
      "Color": "red",
      "CmdStr": "php artisan serve",
      "ExecPath": "."
    },
    {
      "Label": "Vue",
      "Color": "green",
      "CmdStr": "npm run dev",
      "ExecPath": "."
    },
    {
      "Label": "Open VsCode",
      "Color": "blue",
      "CmdStr": "code .",
      "ExecPath": "."
    }
  ]
}
```

### 🔀 Step 4: Run Commands

While in your project directory, simply run:

```bash
xrun
```

### 👇 Output Example

```bash
  _____             _____ _   
 |  __ \           |_   _| |  
 | |__) |   _ _ __   | | | |_
 |  _  / | | | '_ \  | | | __|
 | | \ \ |_| | | | |_| |_| |_
 |_|  \_\__,_|_| |_|_____|\__|
                              
                              
Configured Commands:
Laravel ->   php artisan serve
Vue ->   npm run dev
Open VsCode ->   code .
========================
Running command: Laravel
Running command: Vue
Running command: Open VsCode
```

---

## 🎉 Enjoy Productivity Boost!

With `xrun`, managing your development environment has never been easier! Happy coding! 🚀

---

## 🛠 Feedback and Contributions

Feel free to open issues or submit pull requests. Let’s make `xrun` better together! ✨


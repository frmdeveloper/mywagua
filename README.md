<div align="center">
<img src="assets/bukankahinimy.jpg" width="30%">

# 🚀 My WA Gua

### *WhatsApp Automation Made Simple Piw Piw from Whatsmeow*
---
</div>

## ✨ What is My WA Gua?

Anu, this is a project gabut for my self njir. Don't use this if you are not bored!!

### 🎯 Why Choose This?

- **🔥 Unique** - even though there is already WhatsApp automation available for NodeJS like Baileys, it shoots Whatsmeow instead
---

## 🌟 Features

### Explore for yourself 👍
___
## 🚀 Getting Started

### Prerequisites

- Nodejs 20+
- Go 1.19+ (for building whatsmeow)

### Installation

```bash
npm install mywagua@github:frmdeveloper/mywagua
```

### Quick Start

```javascript
import { makeWASocket } from "mywagua"
//const { makeWASocket } = require("mywagua")

const conn = makeWASocket({
    DbPath: "anu.db",
    Logger: { Database: "DEBUG", Client: "DEBUG", Color: true }
})
console.log(conn)
if (!conn.Store().ID) {
    conn.Connect()
    console.log(conn.PairPhone("62xxx"))
    setInterval(() => {
        if (!conn.IsLoggedIn()) {
            conn.Disconnect()
        }
    }, 60000)
} else {
    conn.Connect()
}
conn.Event(async({type, evt}) => {
    //...
})
```
---

<div align="center">

**Made with 💩**

*If this project helped you, please consider giving it a ⭐ on GitHub!*

</div>

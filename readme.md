# <img src="https://cliply.co/wp-content/uploads/2021/08/372108630_DISCORD_LOGO_400.gif" alt="Logo" width="30"/> DAB (Discord As Database)


# Discord Drive

**DAB** is an API built in Go that allows you to upload and manage files on Discord. This API provides features for segmenting large files, sending them to dedicated channels on Discord, and easily retrieving them.

## Table of Contents

1. [🔀Endpoints](#endpoints)
   - [POST /upload](#post-upload)
   - [GET /list](#get-list)
   - [GET /get](#get-get)
2. [Installation](#installation)

## 🔀Endpoints

### POST /upload 📂

- **Description** :  This route allows you to upload files. Files are split into 10 MB segments and sent to a dedicated channel on Discord. 
- **Réponse** : 
  - `✅200 OK`: The file has been uploaded and segments sent successfully. 🎉
  - `🙅‍♂️400 Bad Request`: Error receiving the file.
  - `❌500 Internal Server Error`:  Error sending segments or other internal issues.

### GET /list 📄

- **Description** : This endpoint lists all files currently stored on Discord.
- **Réponse** : 
  - `✅200 OK`: Returns a list of available files.
  - `❌500 Internal Server Error`: Error retrieving channels.

### GET /get 📥

- **Description** : This endpoint retrieves a file by providing its name. It gathers all segments of the file and serves them for download.
- **Paramètres** : 
  - `filename`: The name of the file to retrieve.
- **Réponse** : 
  - `✅200 OK`: Returns the file for download.
  - `🙅‍♂️400 Bad Request`: If the file does not exist.
  - `❌500 Internal Server Error`: Error reconstructing the file.

## Installation

To install the application, ensure you have Go installed on your system. Then, clone the repository and install the dependencies:

```bash
git clone https://github.com/pepedinho/DAB
cd discord_drive
go run main.go
```

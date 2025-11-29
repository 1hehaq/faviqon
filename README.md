<div align="center">
  <!-- <img src="https://github.com/user-attachments/assets/e87dd5b5-3135-4fed-8a9f-dc75671b85a6" alt="shef" width="955"> -->
</div>

<br>
<br>
<br>

> [!NOTE] 
> **faviqon is a minimal tool for computing favicon hashes from URLs and generating shodan dorks for further reconnaissance.**

<br>

- <sub> **computes murmur3 hashes from favicon.ico files** </sub>
- <sub> **generates shodan search queries for favicon hash matching** </sub>
- <sub> **clean and pipe friendly output** </sub>

<br>
<br>

<h4>Installation</h4>

```bash
go install github.com/haq/faviqon@latest
```

<br>
<br>

<h4>Flags</h4>

<pre>
  -shodan : show shodan dorks only
  -v      : show version
  -h      : show help message
</pre>

<br>
<br>

<h4>Example Commands</h4>

```bash
# compute hashes from URLs
cat https://github.com/ | faviqon
```

<br>

```bash
# generate shodan dorks for the hashes
cat urls.txt | faviqon -shodan
```

<br>

```bash
# pipe to other tools for further processing
shef -q $(echo https://github.com/ | ./faviqon -shodan)
```

<br>
<br>

- **If you see no results or errors**
  - <sub> **verify your URLs are accessible** </sub>
  - <sub> **check your internet connection** (necessary to fetch /favicon.ico) </sub>
  - <sub> **use `-h` for guidance** </sub>

<br>
<br>

> [!CAUTION] 
> **never use `faviqon` for any illegal activities, I'm not responsible for your deeds with it. Do for justice.**

<br>
<br>
<br>

<h6 align="center">kindly for hackers</h6>

<div align="center">
  <a href="https://github.com/1hehaq"><img src="https://img.icons8.com/material-outlined/20/808080/github.png" alt="GitHub"></a>
  <a href="https://twitter.com/1hehaq"><img src="https://img.icons8.com/material-outlined/20/808080/twitter.png" alt="X"></a>
</div>

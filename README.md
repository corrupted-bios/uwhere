# UWhere

UWhere (URL Where) is a tool that follows URLs given through the pipe (stdin) and outputs the final URL redirected.

*That project is mine. You'll see that I forked it from another repo, that other repo is from another account that I have.
## Install

```
▶ go install github.com/davidwkirsch/uwhere@latest
```

## Basic Usage

uwhere accepts line-delimited domains on `stdin`:

```
▶ cat recon/example/domains.txt
https://example.com
https://sub1.example.com
https://sub1.example.com
https://example.edu
https://example.net
▶ cat recon/example/domains.txt | uwhere
https://example.com/index
https://example.com/
https://example.com/
https://example.edu/auth/login
https://account.example.net
```

## Concurrency

You can set the concurrency level with the `-c` flag:

```
▶ cat domains.txt | uwhere -c 50
```

## Timeout

You can change the timeout by using the `-t` flag and specifying a timeout in milliseconds:

```
▶ cat domains.txt | uwhere -t 20000
```


## Why does it exist?

I developed the tool to run on recon of a target, since after grabing multiple domains and subdomains, some of them redirect to the same url (e.g. "sub1.example.com" and "sub2.example.com" redirect to "example.com"). So, in order to eliminate them, that tool will follow the redirect and give you the final url. That gives you the oportunity to get a better result when running a content grabber like fff.

The code was ~~mostly stolen~~ heavily based on the tool [httprobe](https://github.com/tomnomnom/httprobe) from [@tomnomnom](https://github.com/tomnomnom), as the idea. After watching a video on his workflow for reconnaissance, I noticed that after using httprobe he also would use [fff](https://github.com/tomnomnom/fff) (another of his tools), which grabs the head and body of the page given, but if that page is going to redirect you to another directory (e.g. "api.example.com" redirects to "api.example.com/auth/login"), the content grabbed by fff may not be of any value. Also, that gives a quick way to see which url's send you to a login page or something else.

It is my first tool developed ever. I've never seen golang before and I did it in a few hours, so it works, but I bet there's something that can be improved. I'm open to feedbacks and sugestions.

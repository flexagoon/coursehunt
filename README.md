# CourseHunt

This is a website that lets you search for courses across six different platforms.

A school project that I crammed right before the deadline with zero motivation.
If you discover some of the awful things I did in the code, just remember this.

## FAQ

### What sources are supported?

- Udemy
- Stepik
- Coursera
- edX
- Skillbox
- ALISON

### Why does the code suck?

This is already answered above:

> A school project that I crammed right before the deadline with zero motivation.

### What stack does this use?

The Primeagen stack (Go, HTMX, Tailwind)

### Why is it so slow?

Because I spent 0 seconds optimizing it and it makes like 10 API requests for
each query with zero caching.

### How do you access the course data?

Most of the supported platforms (with the exception of Udemy) don't have a public
API, so I just looked at the requests they make using Firefox devtools and tried
to reverse engineer that mess.

### You have a Udemy API key in your code!

I know.

### Do you have any plans to maintain this?

No.

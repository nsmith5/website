FROM scratch
ADD public ./public/
ADD server .
EXPOSE 3000
CMD ["/server", "-dir", "public"]
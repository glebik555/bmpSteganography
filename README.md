# bmpSteganography
## The purpose of the algorithm
Invisible to a third person, the transmission of a message is carried out by sequentially inserting the green pixel component of the next message position into the last bit.
## Algorithm
In order for the receiving side to be able to decrypt the message from the data stream, they need to know **how many components of the green color of the pixel** need to pull out the last bit.  

The sending side calculates the number of pixels of the image, calculates a power of two that is close, but less than the number of pixels. (~~*__less is more__*~~ :alien:)  

The calculated number is **the number of pixels in the green component of which we will "hide" the length of the message**.   

_On the receiving side, the algorithm is similar._

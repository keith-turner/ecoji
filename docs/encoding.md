# Encoji encoding standard

Encoji maps input data into 1024 Unicode emojis (plus 5 padding emojis).  Ten
bits are needed to represent 1024. Ecoji reads 5 bytes at a time because this
is 40 bits which is a multiple of 10.  For each 5 bytes read, 4 emojis are
output.  When less than 5 bytes are available, special padding emojis are
output.  In [mapping.go](../mapping.go) the 1024 emojis and the padding emojis
are defined.  These same emojis should be used in other languages.

Below is some pseudo code for translating bytes to emojis.  Also see [encode.go](../encode.go).

```java

Input input; //input data, read 5 bytes at a time from
Output output; // where unicode emojis are written to
byte data[5]; //buffer that bytes are read into
int numRead;

//assumed this reads maximum available data up to five bytes
while ((numRead = input.read(data)) > 0) {
   for(int i = numRead; i < 5; i++) {
     //zero out unread data
     data[i] = 0;
   }

   switch (numRead) {
      case 1:
        output.writeUnicode(emojis[data[0]<<2 | data[1]>>6]);
        output.writeUnicode(padding);
        output.writeUnicode(padding);
        output.writeUnicode(padding);
        break;
      case 2:
        output.writeUnicode(emojis[data[0]<<2 | data[1]>>6]);
        output.writeUnicode(emojis[(data[1] & 0x3f)<<4 | data[2]>>4]);
        output.writeUnicode(padding);
        output.writeUnicode(padding);
        break;
      case 3:
        output.writeUnicode(emojis[data[0]<<2 | data[1]>>6]);
        output.writeUnicode(emojis[(data[1] & 0x3f)<<4 | data[2]>>4]);
        output.writeUnicode(emojis[(data[2] & 0x0f)<<6 | data[3]>>2]);
        output.writeUnicode(padding);
        break;
      case 4:
        output.writeUnicode(emojis[data[0]<<2 | data[1]>>6]);
        output.writeUnicode(emojis[(data[1] & 0x3f)<<4 | data[2]>>4]);
        output.writeUnicode(emojis[(data[2] & 0x0f)<<6 | data[3]>>2]);
        
        //look at last two bits of 4th byte to determine padding to use
        switch (data[3] & 0x03) {
           case 0:
             output.writeUnicode(padding40);
             break;
           case 1:
             output.writeUnicode(padding41);
             break;
           case 2:
             output.writeUnicode(padding42);
             break;
           case 3:
             output.writeUnicode(padding43);
             break;
        }
        break;

      case 5:
        // use 8 bits from 1st byte and 2 bits from 2nd byte to lookup emoji
        output.writeUnicode(emojis[data[0]<<2 | data[1]>>6]);
        // use 6 bits from 2nd byte and 4 bits from 3rd byte to lookup emoji
        output.writeUnicode(emojis[(data[1] & 0x3f)<<4 | data[2]>>4]);
        // use 4 bits from 3rd byte and 6 bits from 4th byte to lookup emoji
        output.writeUnicode(emojis[(data[2] & 0x0f)<<6 | data[3]>>2]);
        //user 2 bits from 4th byte and 8 bits from 5th byte to lookup emoji
        output.writeUnicode(emojis[(data[3] & 0x03)<<8 | data[4]]);
        break;
   }
}

```
  
For decoding, see [decode.go](../decode.go).  The code needs to be cleaned up, it was written while learning Go.

 

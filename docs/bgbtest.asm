BOOT:0000 31 FE FF         ld   sp,FFFE
BOOT:0003 AF               xor  a
BOOT:0004 21 FF 9F         ld   hl,9FFF
BOOT:0007 32               ldd  (hl),a
BOOT:0008 CB 7C            bit  7,h
BOOT:000A 20 FB            jr   nz,0007
BOOT:000C 21 26 FF         ld   hl,FF26
BOOT:000F 0E 11            ld   c,11
BOOT:0011 3E 80            ld   a,80
BOOT:0013 32               ldd  (hl),a
BOOT:0014 E2               ld   (ff00+c),a
BOOT:0015 0C               inc  c
BOOT:0016 3E F3            ld   a,F3
BOOT:0018 E2               ld   (ff00+c),a
BOOT:0019 32               ldd  (hl),a
BOOT:001A 3E 77            ld   a,77
BOOT:001C 77               ld   (hl),a
BOOT:001D 3E FC            ld   a,FC
BOOT:001F E0 47            ld   (ff00+47),a
BOOT:0021 21 04 01         ld   hl,0104
BOOT:0024 E5               push hl
BOOT:0025 11 CB 00         ld   de,00CB
BOOT:0028 1A               ld   a,(de)
BOOT:0029 13               inc  de
BOOT:002A BE               cp   (hl)
BOOT:002B 20 6B            jr   nz,0098
BOOT:002D 23               inc  hl
BOOT:002E 7D               ld   a,l
BOOT:002F FE 34            cp   a,34
BOOT:0031 20 F5            jr   nz,0028
BOOT:0033 06 19            ld   b,19
BOOT:0035 78               ld   a,b
BOOT:0036 86               add  (hl)
BOOT:0037 23               inc  hl
BOOT:0038 05               dec  b
BOOT:0039 20 FB            jr   nz,0036
BOOT:003B 86               add  (hl)
BOOT:003C 20 5A            jr   nz,0098
BOOT:003E D1               pop  de
BOOT:003F 21 10 80         ld   hl,8010
BOOT:0042 1A               ld   a,(de)
BOOT:0043 CD A9 00         call 00A9
BOOT:0046 CD AA 00         call 00AA
BOOT:0049 13               inc  de
BOOT:004A 7B               ld   a,e
BOOT:004B FE 34            cp   a,34
BOOT:004D 20 F3            jr   nz,0042
BOOT:004F 3E 18            ld   a,18
BOOT:0051 21 2F 99         ld   hl,992F
BOOT:0054 0E 0C            ld   c,0C
BOOT:0056 32               ldd  (hl),a
BOOT:0057 3D               dec  a
BOOT:0058 28 09            jr   z,0063
BOOT:005A 0D               dec  c
BOOT:005B 20 F9            jr   nz,0056
BOOT:005D 11 EC FF         ld   de,FFEC
BOOT:0060 19               add  hl,de
BOOT:0061 18 F1            jr   0054
BOOT:0063 67               ld   h,a
BOOT:0064 3E 64            ld   a,64
BOOT:0066 57               ld   d,a
BOOT:0067 E0 42            ld   (ff00+42),a
BOOT:0069 3E 91            ld   a,91
BOOT:006B E0 40            ld   (ff00+40),a
BOOT:006D 04               inc  b
BOOT:006E 1E 02            ld   e,02
BOOT:0070 CD BC 00         call 00BC
BOOT:0073 0E 13            ld   c,13
BOOT:0075 24               inc  h
BOOT:0076 7C               ld   a,h
BOOT:0077 1E 83            ld   e,83
BOOT:0079 FE 62            cp   a,62
BOOT:007B 28 06            jr   z,0083
BOOT:007D 1E C1            ld   e,C1
BOOT:007F FE 64            cp   a,64
BOOT:0081 20 06            jr   nz,0089
BOOT:0083 7B               ld   a,e
BOOT:0084 E2               ld   (ff00+c),a
BOOT:0085 0C               inc  c
BOOT:0086 3E 87            ld   a,87
BOOT:0088 E2               ld   (ff00+c),a
BOOT:0089 F0 42            ld   a,(ff00+42)
BOOT:008B 90               sub  b
BOOT:008C E0 42            ld   (ff00+42),a
BOOT:008E 15               dec  d
BOOT:008F 20 DD            jr   nz,006E
BOOT:0091 05               dec  b
BOOT:0092 20 69            jr   nz,00FD
BOOT:0094 16 20            ld   d,20
BOOT:0096 18 D6            jr   006E
BOOT:0098 3E 91            ld   a,91
BOOT:009A E0 40            ld   (ff00+40),a
BOOT:009C 1E 14            ld   e,14
BOOT:009E CD BC 00         call 00BC
BOOT:00A1 F0 47            ld   a,(ff00+47)
BOOT:00A3 EE FF            xor  a,FF
BOOT:00A5 E0 47            ld   (ff00+47),a
BOOT:00A7 18 F3            jr   009C
BOOT:00A9 4F               ld   c,a
BOOT:00AA 06 04            ld   b,04
BOOT:00AC C5               push bc
BOOT:00AD CB 11            rl   c
BOOT:00AF 17               rla  
BOOT:00B0 C1               pop  bc
BOOT:00B1 CB 11            rl   c
BOOT:00B3 17               rla  
BOOT:00B4 05               dec  b
BOOT:00B5 20 F5            jr   nz,00AC
BOOT:00B7 22               ldi  (hl),a
BOOT:00B8 23               inc  hl
BOOT:00B9 22               ldi  (hl),a
BOOT:00BA 23               inc  hl
BOOT:00BB C9               ret  
BOOT:00BC 0E 0C            ld   c,0C
BOOT:00BE F0 44            ld   a,(ff00+44)
BOOT:00C0 FE 90            cp   a,90
BOOT:00C2 20 FA            jr   nz,00BE
BOOT:00C4 0D               dec  c
BOOT:00C5 20 F7            jr   nz,00BE
BOOT:00C7 1D               dec  e
BOOT:00C8 20 F2            jr   nz,00BC
BOOT:00CA C9               ret  
BOOT:00CB CE ED            adc  a,ED
BOOT:00CD 66               ld   h,(hl)
BOOT:00CE 66               ld   h,(hl)
BOOT:00CF CC 0D 00         call z,000D
BOOT:00D2 0B               dec  bc
BOOT:00D3 03               inc  bc
BOOT:00D4 73               ld   (hl),e
BOOT:00D5 00               nop  
BOOT:00D6 83               add  e
BOOT:00D7 00               nop  
BOOT:00D8 0C               inc  c
BOOT:00D9 00               nop  
BOOT:00DA 0D               dec  c
BOOT:00DB 00               nop  
BOOT:00DC 08 11 1F         ld   (1F11),sp
BOOT:00DF 88               adc  b
BOOT:00E0 89               adc  c
BOOT:00E1 00               nop  
BOOT:00E2 0E DC            ld   c,DC
BOOT:00E4 CC 6E E6         call z,E66E
BOOT:00E7 DD               undefined opcode
BOOT:00E8 DD               undefined opcode
BOOT:00E9 D9               reti 
BOOT:00EA 99               sbc  c
BOOT:00EB BB               cp   e
BOOT:00EC BB               cp   e
BOOT:00ED 67               ld   h,a
BOOT:00EE 63               ld   h,e
BOOT:00EF 6E               ld   l,(hl)
BOOT:00F0 0E EC            ld   c,EC
BOOT:00F2 CC DD DC         call z,DCDD
BOOT:00F5 99               sbc  c
BOOT:00F6 9F               sbc  a
BOOT:00F7 BB               cp   e
BOOT:00F8 B9               cp   c
BOOT:00F9 33               inc  sp
BOOT:00FA 3E FF            ld   a,FF
BOOT:00FC FF               rst  38
BOOT:00FD 3C               inc  a
BOOT:00FE E0 50            ld   (ff00+50),a

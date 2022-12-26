'****************************************************************
'*  Name    : UNTITLED.BAS                                      *
'*  Author  : [select VIEW...EDITOR OPTIONS]                    *
'*  Notice  : Copyright (c) 2013 [select VIEW...EDITOR OPTIONS] *
'*          : All Rights Reserved                               *
'*  Date    : 21.11.2013                                        *
'*  Version : 1.0                                               *
'*  Notes   :                                                   *
'*          :                                                   *
'****************************************************************
Symbol LCD_EXT_SCL                  PORTD.0           ' SCL connected to this pin
Symbol LCD_EXT_SDA                  PORTD.1           ' SDA connected to this pin
'''------------------- LCD VARIABLE --------------------------------------------------
Symbol LCD_EXT_1_READ_ADR           %01001111    '  PCF8574   A2=1 A1=1 A0=1
Symbol LCD_EXT_1_WRITE_ADR          %01001110    '  PCF8574   A2=1 A1=1 A0=1
Symbol LCD_EXT_2_READ_ADR           %01001101    '  PCF8574   A2=1 A1=1 A0=0
Symbol LCD_EXT_2_WRITE_ADR          %01001100    '  PCF8574   A2=1 A1=1 A0=0
'Symbol LCD_READ1                    %01111101    '  PCF8574A  A2=1 A1=1 A0=0
'Symbol LCD_WRITE1                   %01111100    '  PCF8574A  A2=1 A1=1 A0=0
Symbol LCD_EXT_LINE_1               $00
Symbol LCD_EXT_LINE_2               $40
Symbol LCD_EXT_LINE_3               $14
Symbol LCD_EXT_LINE_4               $54
Dim LCD_BYTE_LOW_NIBBLE             As Byte
Dim LCD_BYTE_HIGH_NIBBLE            As Byte
Dim LCD_EXT_1_I2C                   As Bit
Dim LCD_EXT_2_I2C                   As Bit
Dim LCD_Light                       As Bit
Dim DDRAM_Flag                      As Bit
Dim LCD_EXT_BYTE                    As Byte
Dim LCD_EXT_PnP                     As Byte
Symbol LCD_EXT_RS                   LCD_EXT_BYTE.0
Symbol LCD_EXT_RW                   LCD_EXT_BYTE.1
Symbol LCD_EXT_E                    LCD_EXT_BYTE.2
Symbol LCD_EXT_BackLight            LCD_EXT_BYTE.3
Symbol LCD_EXT_D4                   LCD_EXT_BYTE.4
Symbol LCD_EXT_D5                   LCD_EXT_BYTE.5
Symbol LCD_EXT_D6                   LCD_EXT_BYTE.6
Symbol LCD_EXT_D7                   LCD_EXT_BYTE.7
'LCD_EXT_BYTE = %00111100  bits= D7,D6,D5,D4,LCD_LED,E,RW,RS

GoTo Over_LCD_EXT_I2c_Sub
'-------------------------- Initialize LCD routine -------------------------------------------------------------------------
LCD_EXT_I2C_Init_Sub:
					LCD_EXT_BYTE = %00111100 : GoSub LCD_EXT_I2C_Send_Command_Sub : DelayMS 5     'DB5=1 DB4=1
                    LCD_EXT_BYTE = %00111100 : GoSub LCD_EXT_I2C_Send_Command_Sub : DelayUS 120   'DB5=1 DB4=1
                    LCD_EXT_BYTE = %00111100 : GoSub LCD_EXT_I2C_Send_Command_Sub : DelayUS 120   'DB5=1 DB4=1
                    LCD_EXT_BYTE = %00101100 : GoSub LCD_EXT_I2C_Send_Command_Sub : DelayUS 50    'DB5=1
                    '-----------------------------------------------------------------------
                    LCD_EXT_BYTE = %00101100 : GoSub LCD_EXT_I2C_Send_Command_Sub                 'DB5=1  number of display lines and character font
                    LCD_EXT_BYTE = %10001100 : GoSub LCD_EXT_I2C_Send_Command_Sub : DelayUS 50    'DB7=1  N=1 F=0
                    '-----------------------------------------------------------------------
                    LCD_EXT_BYTE = %00001100 : GoSub LCD_EXT_I2C_Send_Command_Sub                 ' Display on
                    LCD_EXT_BYTE = %11001100 : GoSub LCD_EXT_I2C_Send_Command_Sub : DelayUS 50    ' DB7=1 'DB6=1
                    '-----------------------------------------------------------------------
                    LCD_EXT_BYTE = %00001100 : GoSub LCD_EXT_I2C_Send_Command_Sub
                    LCD_EXT_BYTE = %00011100 : GoSub LCD_EXT_I2C_Send_Command_Sub : DelayUS 2000  ' Clear screen
                    '-----------------------------------------------------------------------
                    LCD_EXT_BYTE = %00001100 : GoSub LCD_EXT_I2C_Send_Command_Sub
                    LCD_EXT_BYTE = %01101100 : GoSub LCD_EXT_I2C_Send_Command_Sub : DelayUS 50    ' Entry mode set
                    For LCD_SGRAM_SYMBOL = 1 To 7
                        GoSub LCD_EXT_I2C_Write_CGRAM_Sub
                    Next                    
                    Return
'-------------------------- Write bitmap to CGRAM routine -------------------------------------------------------------------------
LCD_EXT_I2C_Write_CGRAM_Sub:
                    LCD_DATA_BYTE = (LCD_SGRAM_SYMBOL * 8) + $40                           ' Calculate CGRAM symbol address $40-$78
                    GoSub LCD_EXT_I2C_COMMAND_Sub                                          ' Set CGRAM pointer
'                    For LCD_Index = LCD_DATA_BYTE - $48  To LCD_DATA_BYTE - $40            ' Set pointer to LCD_SGRAM_STRING
                    For LCD_Index = LCD_SGRAM_SYMBOL * 8  To (LCD_SGRAM_SYMBOL * 8) + 8    ' Set pointer to LCD_SGRAM_STRING
                        LCD_DATA_BYTE = LCD_SGRAM_STRING[LCD_Index]
                        GoSub LCD_EXT_I2C_DATA_Sub
                    Next LCD_Index
                    Return
'-------------------------- Write data to DDRAM routine -------------------------------------------------------------------------
'LCD_EXT_I2C_Write_DDRAM_Sub:
LCD_EXT_I2C_PRINT_Sub:
                    ''' SET LINE
                    Select LCD_LINE
                      Case 1
                         LCD_DATA_BYTE = LCD_ROW + LCD_EXT_LINE_1
                      Case 2
                         LCD_DATA_BYTE = LCD_ROW + LCD_EXT_LINE_2
                      Case 3
                         LCD_DATA_BYTE = LCD_ROW + LCD_EXT_LINE_3
                      Case 4
                         LCD_DATA_BYTE = LCD_ROW + LCD_EXT_LINE_4
                    EndSelect                
                    
                    ''' SET ROW
                    LCD_DATA_BYTE = LCD_DATA_BYTE - 1
                    Set DDRAM_Flag : GoSub LCD_EXT_I2C_COMMAND_Sub                           ' Write to DDRAM Flag
                                                                  
                    GoSub Translate
                    
                    END_STRING = Len (STRING_OUT) - 1
                    For LCD_Index = 0 To END_STRING
                        LCD_DATA_BYTE = STRING_OUT[LCD_Index]
                        GoSub LCD_EXT_I2C_DATA_Sub
                    Next LCD_Index
                    Return
'-------------------------- Calculate Command routine --------------------------------------------------------------------------
LCD_EXT_I2C_COMMAND_Sub:
                    LCD_BYTE_LOW_NIBBLE = %00001100                                     ' Calc High Nibble of Command
                    LCD_BYTE_HIGH_NIBBLE = LCD_DATA_BYTE & %11110000
                    LCD_EXT_BYTE = LCD_BYTE_HIGH_NIBBLE | LCD_BYTE_LOW_NIBBLE 
                    If DDRAM_Flag = 1 Then 
                        Set LCD_EXT_D7                                                      ' Set DDRAM pointer
                    Else 
                        Set LCD_EXT_D6                                                      ' Set CGRAM pointer       
                    EndIf
                    GoSub LCD_EXT_I2C_Send_Command_Sub                                      ' Send High Nibble of Command 
                    '-------------------------------------------------------------------------------------------------
                    LCD_BYTE_HIGH_NIBBLE = LCD_DATA_BYTE << 4                           ' Calc Low  Nibble of Command
                    LCD_BYTE_HIGH_NIBBLE = LCD_BYTE_HIGH_NIBBLE & %11110000
                    LCD_EXT_BYTE = LCD_BYTE_HIGH_NIBBLE | LCD_BYTE_LOW_NIBBLE                 
                    GoSub LCD_EXT_I2C_Send_Command_Sub : Clear DDRAM_Flag                   ' Send Low Nibble of Command 
                    Return
'-------------------------- Send Command routine -------------------------------------------------------------------------
LCD_EXT_I2C_Send_Command_Sub:
                    GoSub LCD_EXT_I2C_Send_Strobe_Sub          ' Send High Nibble of Command Strobe 1
                    Clear LCD_EXT_E : Clear LCD_EXT_RS : GoSub LCD_EXT_I2C_Send_Strobe_Sub          ' Send Low  Nibble of Command Strobe 2
                    Return
'-------------------------- Calculate Data routine -------------------------------------------------------------------------
LCD_EXT_I2C_DATA_Sub:
                    LCD_BYTE_LOW_NIBBLE = %00001101                                     ' Calc High Nibble of Data
                    LCD_BYTE_HIGH_NIBBLE = LCD_DATA_BYTE & %11110000
                    LCD_EXT_BYTE = LCD_BYTE_HIGH_NIBBLE | LCD_BYTE_LOW_NIBBLE               ' Write data to CGRAM 
                    GoSub LCD_EXT_I2C_Send_Data_Sub                                         ' Send High Nibble of Data
                    '-------------------------------------------------------------------------------------------------
                    LCD_BYTE_HIGH_NIBBLE = LCD_DATA_BYTE << 4
                    LCD_BYTE_HIGH_NIBBLE = LCD_BYTE_HIGH_NIBBLE & %11110000             ' Calc Low Nibble of Data
                    LCD_EXT_BYTE = LCD_BYTE_HIGH_NIBBLE | LCD_BYTE_LOW_NIBBLE               ' Write data to CGRAM 
                    GoSub LCD_EXT_I2C_Send_Data_Sub                                         ' Send Low Nibble of Data
                    Return
'-------------------------- Send Data routine -------------------------------------------------------------------------
LCD_EXT_I2C_Send_Data_Sub:
                    GoSub LCD_EXT_I2C_Send_Strobe_Sub           ' Send High Nibble of Data Strobe 1
                    Clear LCD_EXT_E : GoSub LCD_EXT_I2C_Send_Strobe_Sub                         ' Send Low  Nibble of Data Strobe 2
                    Return
'-------------------------- Clear screen routine -------------------------------------------------------------------------
LCD_EXT_I2C_CLS_Sub:
                    GoSub LCD_EXT_I2C_Init_Sub
'                    LCD_EXT_BYTE = %00001100 : GoSub LCD_EXT_I2C_Send_Command_Sub
'                    LCD_EXT_BYTE = %00011100 : GoSub LCD_EXT_I2C_Send_Command_Sub : DelayUS 1800  ' Clear screen
                    Return
'-------------------------- LCD backlight On/Off routine -------------------------------------------------------------------------
'LCD_EXT_On_Sub:
'                    Set LCD_Light : GoSub LCD_EXT_I2C_Send_Strobe_Sub 'GoSub LCD_EXT_I2C_Send_Command_Sub
'                    Return
'LCD_EXT_Off_Sub:
'                    Clear LCD_Light : GoSub LCD_EXT_I2C_Send_Strobe_Sub 'GoSub LCD_EXT_I2C_Send_Command_Sub
'                    Return
'-------------------------- Send I2C routine ------------------------------------------------------------------------
LCD_EXT_I2C_Send_Strobe_Sub:
                    If LCD_EXT_1_I2C = YES And LCD_NUM = LCD_Attach Then
'                        LCD_EXT_BackLight = 1 
                        I2COut LCD_EXT_SDA, LCD_EXT_SCL, LCD_EXT_1_WRITE_ADR,[LCD_EXT_BYTE]       ' Write Software I2C
                    EndIf                           

'                    If LCD_EXT_2_I2C = YES And LCD_NUM = "B" Then
''                        LCD_EXT_BackLight = LCD_Light 
'                        I2COut LCD_EXT_SDA, LCD_EXT_SCL, LCD_EXT_2_WRITE_ADR,[LCD_EXT_BYTE]
'                    EndIf 
                    Return
'-------------------------- I2C Plug-n-Play routine ------------------------------------------------------------------------
LCD_EXT_I2C_PNP_Sub:
                    I2COut LCD_EXT_SDA, LCD_EXT_SCL, LCD_EXT_1_WRITE_ADR, [%01010101]
                    I2CIn  LCD_EXT_SDA, LCD_EXT_SCL, LCD_EXT_1_READ_ADR, [LCD_EXT_PnP]
                    If LCD_EXT_PnP = %01010101 Then 
                       LCD_EXT_1_I2C = YES
                       If LCD_Attach = "A" Then LCD_A_Config = LCD_Config : LCD_Connect_A = "L" : EndIf 
                       If LCD_Attach = "B" Then LCD_B_Config = LCD_Config : LCD_Connect_B = "L" : EndIf
                    EndIf   
                    Return
'-------------------------- Macros -------------------------------------------------------------------------

'LCD_EXT_Cls  Macro 
'      GoSub LCD_EXT_I2C_CLS_Sub
'Endm

Print_EXT  Macro 
      LCD_NUM = LCD_Attach
      GoSub LCD_EXT_I2C_PRINT_Sub
Endm

'LCD_EXT_Init  Macro 
'      GoSub LCD_EXT_I2C_Init_Sub
'Endm     

'LCD_EXT_On Macro
'      GoSub LCD_EXT_On_Sub
'Endm      

'LCD_EXT_Off Macro
'      GoSub LCD_EXT_Off_Sub
'Endm      
'---------------------------------------------------------------------------------------------------
Over_LCD_EXT_I2c_Sub:

GoSub LCD_EXT_I2C_PNP_Sub
If LCD_EXT_1_I2C = YES Then LCD_NUM = LCD_Attach : GoSub LCD_EXT_I2C_Init_Sub : EndIf
'If LCD_EXT_1_I2C = YES Then LCD_NUM = "B" : GoSub LCD_EXT_I2C_Init_Sub : EndIf
LCD_LINE = 1 : GoSub StartScreen1_Sub


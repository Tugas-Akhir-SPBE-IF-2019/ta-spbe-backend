# this code is for TEXT FINDING - INDIKATOR 9
# input is txt file from preprocess dokumen lama, nama instansi & judul from dokumen lama & baru

import re
import preprocess_dokbaru as dokbaru
import preprocess_doklama as doklama
from .utility import *


def txtreader(filename, lv, keyword):
    # func to search keyword in txt file

    # open txt file
    # additional hacks
    filename = filename.split("/")[-1]
    file = open(f'cleaned_{filename}.txt', 'r')

    idx = 0
    result = []

    # read line by line from txt
    for line in file:

        if (lv == 1):
            for key in keyword:
                reg = f'(?:(audit)\s({key}))'
                if (re.search(reg, line, re.IGNORECASE)):
                    result.append([idx, line])

        else:  # lv == 2
            # cek dgn keyword 'audit'
            for key in keyword:

                # manfaatin regex untuk searching yg not only kata keyword
                reg = f'(?:(audit)\s({key}))'
                if (re.search(reg, line, re.IGNORECASE)):
                    result.append([idx, line])

        idx += 1

    file.close()
    return (result)


# if __name__ == '__main__':
#     pdfparser(sys.argv[1])

def ceklvl(filename):
    list_final = []
    text_final = ''

    lvl1 = ["TIK", "Teknologi Informasi dan Komunikasi"]
    res1 = txtreader(filename, 1, lvl1)

    # cek if keyword lvl1 is not found, then return as empty string
    if (not res1):
        return text_final

    lvl2 = convert_keywords(
        ["Infrastruktur SPBE", "Aplikasi SPBE", "Keamanan SPBE"])
    res2 = txtreader(filename, 2, lvl2)

    if (not res2):
        list_final.append(res1[0][1])
        return clean_text(list_final)

    for el in res2:
        if (el[1] not in list_final):
            list_final.append(el[1])

    text_final = ". ".join(list_final)

    return clean_text(list_final)
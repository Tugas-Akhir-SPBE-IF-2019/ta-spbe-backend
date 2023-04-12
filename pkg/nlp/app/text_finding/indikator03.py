# this code is for TEXT FINDING - INDIKATOR 3
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
    initial_keys = list(keyword)

    # read line by line from txt
    for line in file:

        if (lv == 1):
            if re.search(keyword[0], line, re.IGNORECASE):
                result.append([idx, line])

        elif (lv == 3):
            # Check primary word first:
            # 1. "Manajemen Data"
            # 2. "Data dan Informasi"
            reg1 = f'(?:(manajemen)\s+(data))'
            reg2 = f'(?:(data)\s+(dan)\s+(informasi))'

            if (re.search(reg1, line, re.IGNORECASE) or re.search(reg2, line, re.IGNORECASE)):
                # Use regex to handle multiple whitespaces
                exclude_keywords = ["sistem", "domain", "aplikasi"]

                for key in initial_keys:
                    if (re.search(key, line, re.IGNORECASE) and not exclude_words(exclude_keywords, line)):
                        result.append([idx, line])

                        # Delete found key to prevent redundant search
                        if (key in keyword):
                            keyword.remove(key)

        else:  # lv == 5
            # Check primary word first:
            # 1. Manajemen Data SPBE
            reg = f'(?:(manajemen)\s+(data)\s+(spbe))'
            if (re.search(reg, line, re.IGNORECASE)):
                res = [ele for ele in keyword if (ele in line)]
                if (res):
                    result.append([idx, line])

        idx += 1

    file.close()
    return result

def ceklvl(filename):
    list_final = []

    # Same keyword for lvl1 and lvl2
    lvl1 = convert_keywords(["Manajemen Data"])
    res1 = txtreader(filename, 1, lvl1)

    # cek if keyword lvl1 is not found, then return as empty string
    if (not res1):
        return ''

    # Same keyword for lvl3 and lvl4
    lvl3 = convert_keywords([
        "Arsitektur Data", "Data Induk", "Data Referensi",
        "Basis Data", "Kualitas Data", "Interoperabilitas Data",
        "Interoperabilitas"
    ])
    res3 = txtreader(filename, 3, lvl3)

    # Return result for level 2 if no Level 3 Keywords found
    if (not res3):
        # Add first line found containing "manajemen data" to list_final for level 2
        list_final.append(res1[0][1])
        return clean_text(list_final)

    for el in res3:
        if (el[1] not in list_final):
            list_final.append(el[1])
    
    # Terminate immediately if not all Level 3 Keywords found
    if not(len(lvl3) == 0 or (len(lvl3) == 1 and "Interoperabilitas" in lvl3[0])):
        return clean_text(list_final)

    lvl5 = ["integrasi", "reviu", "diselaraskan", "berpedoman", "pedoman", "perubahan"]
    res5 = txtreader(filename, 5, lvl5)

    for el in res5:
        if (el[1] not in list_final):
            list_final.append(el[1])

    # return text_final
    return clean_text(list_final)
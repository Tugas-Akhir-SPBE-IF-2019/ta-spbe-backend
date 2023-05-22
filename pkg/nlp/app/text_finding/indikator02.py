# this code is for TEXT FINDING - INDIKATOR 2
# input is txt file from preprocess dokumen lama, nama instansi & judul from dokumen lama & baru

import re
from .utility import *


def txtreader(filename, lv, keyword):
    # func to search keyword in txt file

    # open txt file
    # additional hacks
    filename = filename.split("/")[-1]
    file = open(f'cleaned_{filename}.txt', 'r')

    idx = 0
    result = []
    initial_keywords = list(keyword)

    # read line by line from txt
    for line in file:

        if (lv == 1):
            if re.search(keyword[0], line, re.IGNORECASE):
                result.append([idx, line])

        elif (lv == 2):
            # Check primary words first:
            # 1. "Peta Rencana SPBE"
            # 2. "Rencana Strategis"
            reg1 = f'(?:(peta)\s+(rencana)\s+(spbe))'
            reg2 = f'(?:(rencana)\s+(strategis))'

            if (re.search(reg1, line, re.IGNORECASE) or re.search(reg2, line, re.IGNORECASE)):
                # Use regex to handle multiple whitespaces

                for key in initial_keywords:
                    if (re.search(key, line, re.IGNORECASE)):
                        result.append([idx, line])

                        # Delete found key to prevent redundant search
                        if (key in keyword):
                            keyword.remove(key)

        else:  # lv == 4
            # Check primary words first:
            # 1. "Peta Rencana SPBE Nasional"
            # 2. "Rencana Induk SPBE Nasional"
            reg1 = f'(?:(peta)\s+(rencana)\s+(spbe)\s+(nasional))'
            reg2 = f'(?:(rencana)\s+(induk)\s+(spbe)\s+(nasional))'

            if (re.search(reg1, line, re.IGNORECASE) or re.search(reg2, line, re.IGNORECASE)):
                res = [ele for ele in keyword if (ele in line)]
                if (res):
                    result.append([idx, line])

        idx += 1

    file.close()
    return result


def ceklvl(filename):
    list_final = []

    lvl1 = convert_keywords(["Peta Rencana SPBE"])
    res1 = txtreader(filename, 1, lvl1)

    # cek if keyword lvl1 is not found, then return as empty string
    if (not res1):
        return ''

    # Regex keywords to handle multiple whitespaces
    lvl2 = convert_keywords([
        "Tata Kelola SPBE", "Manajemen SPBE",
        "Layanan SPBE", "Infrastruktur SPBE",
        "Aplikasi SPBE", "Keamanan SPBE",
        "Audit Teknologi Informasi dan Komunikasi",
        "Audit TIK"
    ])
    res2 = txtreader(filename, 2, lvl2)

    # Terminate immediately if no Level 2 Keywords found
    if (not res2):
        # Add first line found containing "peta rencana spbe" to list_final for level 2
        list_final.append(res1[0][1])
        return clean_text(list_final)

    for el in res2:
        if (el[1] not in list_final):
            list_final.append(el[1])
    
    # Terminate if not all Level 3 Keywords found
    if not(len(lvl2) == 0 or (len(lvl2) == 1 and "Audit" in lvl2[0])):
        return clean_text(list_final)

    lvl4 = ["integrasi", "reviu", "diselaraskan", "berpedoman", "perubahan"]
    res4 = txtreader(filename, 4, lvl4)

    for el in res4:
        list_final.append(el[1])

    return clean_text(list_final)
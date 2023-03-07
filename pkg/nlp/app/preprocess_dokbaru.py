# this code is for DOKUMEN BUKTI BARU
# the code below will convert it from pdf into readable data
# resulting in txt file where each line can be seen as separate entities, to further be processed for text finding

import sys
from pdfminer.pdfinterp import PDFResourceManager, PDFPageInterpreter
from pdfminer.pdfpage import PDFPage
from pdfminer.converter import XMLConverter, HTMLConverter, TextConverter
from pdfminer.layout import LAParams
import io
import re


def pdfparser(filename):
    path = './src/static/'
    fp = open(path + filename, 'rb')
    rsrcmgr = PDFResourceManager()
    retstr = io.StringIO()
    codec = 'utf-8'
    laparams = LAParams()
    device = TextConverter(rsrcmgr, retstr, codec=codec, laparams=laparams)
    # Create a PDF interpreter object.
    interpreter = PDFPageInterpreter(rsrcmgr, device)
    # Process each page contained in the document.

    page_no_title = 0
    page_no_start = 0
    startingPoint = None
    startfound = False
    fulltext = ''
    judul = ''
    instansi = ''

    for pageNumber, page in enumerate(PDFPage.get_pages(fp)):
        interpreter.process_page(page)
        data = retstr.getvalue()
        retstr.truncate(0)
        # retstr.seek(0)

        if (startfound == False):
            if (pageNumber == page_no_title):
                if (re.search('rancangan', data, re.IGNORECASE) is None):

                    # get nama instansi
                    instansi = re.split('(?i)peraturan', data)[0]

                    # get starting point for judul
                    peraturan = re.split('(?i)(peraturan)', data)[1]
                    data = re.split('(?i)(peraturan)', data)[2]

                    # get judul
                    judul = peraturan + \
                        re.split('(?i)(dengan rahmat)', data)[0]
                    judul = " ".join(judul.split())

                else:

                    # get nama instansi
                    instansi = re.split('(?i)rancangan', data)[0]

                    # get starting point for judul
                    rancangan = re.split('(?i)(rancangan)', data)[1]
                    data = re.split('(?i)(rancangan)', data)[2]

                    # get judul
                    judul = rancangan + \
                        re.split('(?i)(dengan rahmat)', data)[0]
                    judul = " ".join(judul.split())

            # get starting point for isi dokumen bukti
            startingPoint = re.search('(?i)(bab I)', data)

            # if starting point found in current page
            if (startingPoint != None):
                startfound = True
                page_no_start = pageNumber

                # get isi dokumen bukti
                startingPoint = re.split('(?i)(bab I)', data)[1]
                fulltext += startingPoint + re.split('(?i)(bab I)', data)[2]
        else:
            fulltext += data

    # clean data and write to txt
    datacleaner(fulltext, filename)

    return (instansi, judul)


def datacleaner(data, filename):
    codec = 'utf-8'

    # clean whitespace and undefined texts
    data = re.sub(r'\x00+', '\n', data)
    data = re.sub(r'\n{3,}', '\n' * 2, data)
    data = re.sub(r'\n{2,3}', '\n', data)
    data = re.sub(r'(\n\s+)+\n', '\n' * 2, data, flags=re.MULTILINE)

    data = re.sub(r'(\n)+', '', data, flags=re.MULTILINE)

    data = re.split(r'(?:[?!\n]|(?<!\d)\.(?!\d))', data)
    # TODO
    # exclude alphabet numbering from split

    with open(f'cleaned_{filename}.txt', 'w') as file:
        for el in data:

            # TODO
            # how to remove 'empty' line from reading image, table, etc

            # check if not 'empty' line
            if (re.fullmatch('^(?:\s*(atau|dan)?\s?([0-9]+|[a-z]{1,2}))', el) is None):
                # remove leading whitespace and add newline
                # write line to txt
                file.write(f"{el.lstrip()}\n")

# if __name__ == '__main__':
#     pdfparser(sys.argv[1])


# filename = 'F2201-287-Indikator_01~+~Indikator1_Perbup_81_tahun_2021.pdf'
# filename2 = 'Draft Perbup SPBE 2021 Revisi.pdf'

# (instansi, judul) = pdfparser(filename2)

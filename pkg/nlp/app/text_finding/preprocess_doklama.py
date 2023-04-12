# this code is for DOKUMEN BUKTI LAMA
# the code below will get the title

import sys
from pdfminer.pdfinterp import PDFResourceManager, PDFPageInterpreter
from pdfminer.pdfpage import PDFPage
from pdfminer.converter import XMLConverter, HTMLConverter, TextConverter
from pdfminer.layout import LAParams
import io
import re


def pdfparser(data):

    fp = open(data, 'rb')
    rsrcmgr = PDFResourceManager()
    retstr = io.StringIO()
    codec = 'utf-8'
    laparams = LAParams()
    device = TextConverter(rsrcmgr, retstr, codec=codec, laparams=laparams)
    # Create a PDF interpreter object.
    interpreter = PDFPageInterpreter(rsrcmgr, device)
    # Process each page contained in the document.

    page_no_title = 0

    for pageNumber, page in enumerate(PDFPage.get_pages(fp)):
        interpreter.process_page(page)
        data = retstr.getvalue()
        retstr.truncate(0)
        # retstr.seek(0)

        if (pageNumber == page_no_title):
            # get nama instansi
            instansi = re.split('(?i)peraturan', data)[0]
            instansi = stringcleaner(instansi)

            # get starting point for judul
            peraturan = re.split('(?i)(peraturan)', data)[1]
            data = re.split('(?i)(peraturan)', data)[2]

            # get judul
            judul = peraturan + re.split('(?i)(dengan rahmat)', data)[0]
            judul = " ".join(judul.split())
            judul = re.split('(?i)(menimbang)', judul)[0]
            judul = stringcleaner(judul)

        break

    return (instansi, judul)


def stringcleaner(data):
    codec = 'utf-8'

    # clean whitespace and undefined texts
    data = re.sub(r'\x00+', '', data)
    data = re.sub(r'\n{3,}', '', data)
    data = re.sub(r'\n{2,3}', '', data)
    data = re.sub(r'(\n\s+)+\n', '', data, flags=re.MULTILINE)

    data = re.sub(r'(\s)+', ' ', data, flags=re.MULTILINE)
    data = re.sub(r'(\n)+', '', data, flags=re.MULTILINE)

    return data.strip()

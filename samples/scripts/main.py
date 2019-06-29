import string
import re
import hashlib
import shutil
import os
from urllib import request
import nltk

# URL pointing to raw text book to download text and parse
BOOK_URL = "http://www.gutenberg.org/files/2554/2554-0.txt"

# We keep generated paragraph files in this dir first
# Among them selected ones are moved to use directory
TRY_DIR = '../try'

# Includes
# uppercase,lowercase ascii
# digits, punctuation, space
# VALID_CHAR_SET = set(string.printable)
valid_chars_re = re.compile('[^%s]' % string.printable)

# leaves whitespace chars
# Just to make sure to remove
# wierd characters
def remove_nontypable(p):
    return valid_chars_re.sub('', p)

# Decide whether a paragraph should be part of 
# set of samples
def paragraph_approved(p, space_tokenizer):
    words = space_tokenizer.tokenize(p)
    if len(words) < 12 or len(words) > 50: return False
    # TODO add numeric chars % in p check
    # TODO Detect if p is actually a good paragraph. Check nlkt API
    return True

def clear_try_dir():
    if not os.path.exists(TRY_DIR):
        print('Creating /try directory..')
        os.mkdir(TRY_DIR)
    else:
        shutil.rmtree(TRY_DIR)

def save_this_paragraph(p):
    # generate file name. Eg 73824146.txt
    file_name = ('%s/%s.txt' % 
        (TRY_DIR, str(int(hashlib.sha256(p.encode('utf-8')).hexdigest(), 16) % 10**8))
    )
    print('Saving paragraph file %s' % file_name)
    f = open(file_name, 'w', encoding='utf-8')
    f.write(p)
    f.close()
    return 1

def process_paragraph(p, space_tokenizer):
    p = remove_nontypable(p)
    if paragraph_approved(p, space_tokenizer):
        return save_this_paragraph(p)
    return 0

def get_raw_file():
    response = request.urlopen(BOOK_URL)
    raw = response.read().decode('utf8')
    return raw

def main():
    print('Tokenizing...')
    clear_try_dir()

    raw = get_raw_file()
    paragraph_tokenizer = nltk.tokenize.BlanklineTokenizer()
    whitespace_tokenizer = nltk.tokenize.WhitespaceTokenizer()
    paragraphs = paragraph_tokenizer.tokenize(raw)
    generated_count = 0
    for p in paragraphs:
        generated_count += process_paragraph(p, whitespace_tokenizer)

    print ('Generated %d paragraphs.' % generated_count)
    print('Done tokenizing.')

if __name__ == '__main__':
    main()
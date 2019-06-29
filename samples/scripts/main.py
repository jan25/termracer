import string
from urllib import request

import nltk

ALPHABETS = set(string.ascii_lowercase + string.ascii_uppercase)

def _is_word(w):
    if len(set(w).intersection(ALPHABETS)) > 0:
        return True
    return False

# somajo splits even at punctuation marks
# this patch fn tries to put together the punctuation
# with previous word. Example ['world', '.'] => ['world.']
# def _monkey_patch_paragraph(p):
#     p_final = []
#     for w in p:
#         if not _is_word(w) and len(p_final) > 0:
#             p_final[-1] += w
#         else: p_final.append(w)
#     return p_final

# Decide whether a paragraph should be part of 
# set of samples
def _approve_paragraph(p, space_tokenizer):
    # TODO
    words = space_tokenizer.tokenize(p)
    return True

def save_this_paragraph(p):
    pass

def processs_paragraph(p, space_tokenizer):
    print('==============================')
    if _approve_paragraph(p, space_tokenizer):
        save_this_paragraph(p)

def get_raw_file():
    url = "http://www.gutenberg.org/files/2554/2554-0.txt"
    response = request.urlopen(url)
    raw = response.read().decode('utf8')
    return raw

def main():
    print('Tokenizing...')
    
    raw = get_raw_file()
    paragraph_tokenizer = nltk.tokenize.BlanklineTokenizer()
    whitespace_tokenizer = nltk.tokenize.WhitespaceTokenizer()
    paragraphs = paragraph_tokenizer.tokenize(raw)
    for p in paragraphs:
        processs_paragraph(p, whitespace_tokenizer)

    print('Done tokenizing.')

if __name__ == '__main__':
    main()
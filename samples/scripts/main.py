
from somajo import Tokenizer, SentenceSplitter

def processs_paragraph(p):
    pass

def main():
    tokenizer = Tokenizer(language="en")

    # tokenize file to paragraphs
    # assumes paragraphs are seperated by emptylines
    paragraphs = tokenizer.tokenize_file('test.txt')
    for p in paragraphs:
        processs_paragraph(p)

if __name__ == '__main__':
    main()
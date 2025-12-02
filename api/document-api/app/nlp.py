"""
Copyright (c) 2023-2025 RapidaAI
Author: Prashant Srivastav <prashant@rapida.ai>

Licensed under GPL-2.0 with Rapida Additional Terms.
See LICENSE.md for details or contact sales@rapida.ai for commercial use.
"""
import spacy
import nltk
import torch

# Load the SpaCy model
# do not get confuse with faster and accrate
en_core_web_trf = spacy.load("en_core_web_trf")
en_core_web_sm = spacy.load("en_core_web_sm")

torch.set_num_threads(1)


def init_model():
    nltk.download("punkt")

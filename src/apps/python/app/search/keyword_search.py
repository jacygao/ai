import re
from typing import List, Tuple
from rank_bm25 import BM25Okapi
from app.models import Document, SearchResult
from app.utils.data_loader import data_loader

class KeywordSearch:
    def __init__(self):
        self.documents = data_loader.get_documents()
        self.bm25 = None
        self._build_index()
    
    def _build_index(self):
        """Build BM-25 index from documents"""
        if not self.documents:
            return
        
        # Tokenize documents
        tokenized_docs = []
        for doc in self.documents:
            # Simple tokenization - split on whitespace and remove punctuation
            tokens = re.findall(r'\b\w+\b', doc.content.lower())
            tokenized_docs.append(tokens)
        
        # Create BM-25 index
        self.bm25 = BM25Okapi(tokenized_docs)
    
    def search(self, query: str, limit: int = 10) -> List[SearchResult]:
        """Perform keyword search using BM-25"""
        if not self.bm25 or not self.documents:
            return []
        
        # Tokenize query
        query_tokens = re.findall(r'\b\w+\b', query.lower())
        
        if not query_tokens:
            return []
        
        # Get BM-25 scores
        scores = self.bm25.get_scores(query_tokens)
        
        # Create list of (score, document_index) tuples
        doc_scores = [(score, i) for i, score in enumerate(scores) if score > 0]
        
        # Sort by score (descending)
        doc_scores.sort(key=lambda x: x[0], reverse=True)
        
        # Get top results
        results = []
        for score, doc_idx in doc_scores[:limit]:
            doc = self.documents[doc_idx]
            
            # Highlight matching terms in content
            highlighted_content = self._highlight_terms(doc.content, query_tokens)
            
            result = SearchResult(
                id=doc.id,
                title=doc.title,
                content=doc.content,
                category=doc.category,
                tags=doc.tags,
                score=float(score),
                highlighted_content=highlighted_content
            )
            results.append(result)
        
        return results
    
    def _highlight_terms(self, content: str, query_tokens: List[str]) -> str:
        """Highlight matching terms in content"""
        highlighted = content
        
        for token in query_tokens:
            # Create case-insensitive regex pattern
            pattern = re.compile(r'\b' + re.escape(token) + r'\b', re.IGNORECASE)
            highlighted = pattern.sub(f'<mark>{token}</mark>', highlighted)
        
        return highlighted
    
    def get_document_count(self) -> int:
        """Get total number of documents"""
        return len(self.documents) 
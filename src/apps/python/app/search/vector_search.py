import time
from typing import List
import numpy as np
import torch
from transformers import AutoTokenizer, AutoModel
from sklearn.metrics.pairwise import cosine_similarity
from app.models import Document, SearchResult
from app.utils.data_loader import data_loader

class VectorSearch:
    def __init__(self, model_name: str = "sentence-transformers/all-MiniLM-L6-v2"):
        self.documents = data_loader.get_documents()
        self.model = None
        self.tokenizer = None
        self.document_embeddings = None
        self._load_model()
        self._build_embeddings()
    
    def _load_model(self):
        """Load the transformer model and tokenizer"""
        try:
            print("Loading transformer model...")
            self.tokenizer = AutoTokenizer.from_pretrained("sentence-transformers/all-MiniLM-L6-v2")
            self.model = AutoModel.from_pretrained("sentence-transformers/all-MiniLM-L6-v2")
            print("Model loaded successfully!")
        except Exception as e:
            print(f"Error loading model: {e}")
            # Fallback to a simpler model
            try:
                print("Trying fallback model...")
                self.tokenizer = AutoTokenizer.from_pretrained("distilbert-base-uncased")
                self.model = AutoModel.from_pretrained("distilbert-base-uncased")
                print("Fallback model loaded successfully!")
            except Exception as e2:
                print(f"Error loading fallback model: {e2}")
                self.model = None
                self.tokenizer = None
    
    def _get_embeddings(self, texts: List[str]) -> np.ndarray:
        """Generate embeddings for a list of texts"""
        if not self.model or not self.tokenizer:
            return np.array([])
        
        embeddings = []
        
        for text in texts:
            try:
                # Tokenize the text
                inputs = self.tokenizer(
                    text, 
                    return_tensors="pt", 
                    max_length=512, 
                    truncation=True, 
                    padding=True
                )
                
                # Generate embeddings
                with torch.no_grad():
                    outputs = self.model(**inputs)
                    # Use mean pooling of the last hidden state
                    embedding = outputs.last_hidden_state.mean(dim=1).squeeze().numpy()
                    embeddings.append(embedding)
                    
            except Exception as e:
                print(f"Error generating embedding for text: {e}")
                # Use zero vector as fallback
                embeddings.append(np.zeros(self.model.config.hidden_size))
        
        return np.array(embeddings)
    
    def _build_embeddings(self):
        """Build embeddings for all documents"""
        if not self.model or not self.documents:
            self.document_embeddings = None
            return
        
        try:
            print("Building document embeddings...")
            # Create combined text (title + content) for better semantic matching
            texts = [f"{doc.title}. {doc.content}" for doc in self.documents]
            
            # Generate embeddings
            self.document_embeddings = self._get_embeddings(texts)
            print(f"Built embeddings for {len(self.documents)} documents")
        except Exception as e:
            print(f"Error building embeddings: {e}")
            self.document_embeddings = None
    
    def search(self, query: str, limit: int = 10) -> List[SearchResult]:
        """Perform vector search using transformer embeddings"""
        if not self.model or self.document_embeddings is None or not self.documents:
            return []
        
        try:
            # Generate query embedding
            query_embedding = self._get_embeddings([query])
            
            if query_embedding.size == 0:
                return []
            
            # Calculate cosine similarities
            similarities = cosine_similarity(query_embedding, self.document_embeddings)[0]
            
            # Create list of (similarity, document_index) tuples
            doc_scores = [(sim, i) for i, sim in enumerate(similarities)]
            
            # Sort by similarity (descending)
            doc_scores.sort(key=lambda x: x[0], reverse=True)
            
            # Get top results
            results = []
            for similarity, doc_idx in doc_scores[:limit]:
                doc = self.documents[doc_idx]
                
                # For vector search, we don't highlight specific terms
                # but we can show the most relevant parts of the content
                highlighted_content = self._get_relevant_excerpt(doc.content, query)
                
                result = SearchResult(
                    id=doc.id,
                    title=doc.title,
                    content=doc.content,
                    category=doc.category,
                    tags=doc.tags,
                    score=float(similarity),
                    highlighted_content=highlighted_content
                )
                results.append(result)
            
            return results
            
        except Exception as e:
            print(f"Error during vector search: {e}")
            return []
    
    def _get_relevant_excerpt(self, content: str, query: str, max_length: int = 200) -> str:
        """Get a relevant excerpt from the content"""
        # Simple excerpt: take the first part of the content
        if len(content) <= max_length:
            return content
        
        # Try to find a good breaking point
        words = content.split()
        excerpt_words = words[:max_length//5]  # Approximate word count
        
        excerpt = " ".join(excerpt_words)
        if len(excerpt) > max_length:
            excerpt = excerpt[:max_length-3] + "..."
        else:
            excerpt += "..."
        
        return excerpt
    
    def get_document_count(self) -> int:
        """Get total number of documents"""
        return len(self.documents) 
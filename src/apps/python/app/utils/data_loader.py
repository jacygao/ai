import json
import os
from typing import List
from app.models import Document

class DataLoader:
    def __init__(self, data_file: str = "data/sample_documents.json"):
        self.data_file = data_file
        self.documents: List[Document] = []
        self.load_documents()
    
    def load_documents(self) -> None:
        """Load documents from JSON file"""
        try:
            # Get the absolute path to the data file
            current_dir = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
            data_path = os.path.join(current_dir, self.data_file)
            
            with open(data_path, 'r', encoding='utf-8') as f:
                data = json.load(f)
                self.documents = [Document(**doc) for doc in data]
        except FileNotFoundError:
            print(f"Warning: Data file {self.data_file} not found. Using empty document set.")
            self.documents = []
        except Exception as e:
            print(f"Error loading documents: {e}")
            self.documents = []
    
    def get_documents(self) -> List[Document]:
        """Get all loaded documents"""
        return self.documents
    
    def get_document_by_id(self, doc_id: int) -> Document:
        """Get a specific document by ID"""
        for doc in self.documents:
            if doc.id == doc_id:
                return doc
        return None
    
    def get_document_texts(self) -> List[str]:
        """Get list of document texts for indexing"""
        return [doc.content for doc in self.documents]
    
    def get_document_titles(self) -> List[str]:
        """Get list of document titles"""
        return [doc.title for doc in self.documents]
    
    def get_document_ids(self) -> List[int]:
        """Get list of document IDs"""
        return [doc.id for doc in self.documents]
    
    def get_document_count(self) -> int:
        """Get total number of documents"""
        return len(self.documents)

# Global instance
data_loader = DataLoader() 
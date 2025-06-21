from pydantic import BaseModel
from typing import List, Optional
from enum import Enum

class SearchMode(str, Enum):
    KEYWORD = "keyword"
    VECTOR = "vector"

class Document(BaseModel):
    id: int
    title: str
    content: str
    category: str
    tags: List[str]

class SearchRequest(BaseModel):
    query: str
    mode: SearchMode
    limit: Optional[int] = 10

class SearchResult(BaseModel):
    id: int
    title: str
    content: str
    category: str
    tags: List[str]
    score: float
    highlighted_content: Optional[str] = None

class SearchResponse(BaseModel):
    results: List[SearchResult]
    total_results: int
    query: str
    mode: SearchMode
    search_time_ms: float 
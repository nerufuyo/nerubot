"""Whale alert model."""
from dataclasses import dataclass
from datetime import datetime
from typing import Optional, List
from enum import Enum


class AlertType(Enum):
    """Types of whale alerts."""
    LARGE_TRANSFER = "large_transfer"
    DEX_SWAP = "dex_swap"
    EXCHANGE_DEPOSIT = "exchange_deposit"
    EXCHANGE_WITHDRAWAL = "exchange_withdrawal"
    WHALE_ACCUMULATION = "whale_accumulation"
    WHALE_DISTRIBUTION = "whale_distribution"


@dataclass
class WhaleAlert:
    """Class representing a whale alert."""
    
    alert_id: str
    alert_type: AlertType
    token_symbol: str
    token_name: str
    amount: float
    amount_usd: float
    from_address: Optional[str] = None
    to_address: Optional[str] = None
    from_label: Optional[str] = None
    to_label: Optional[str] = None
    transaction_hash: Optional[str] = None
    blockchain: Optional[str] = None
    timestamp: Optional[datetime] = None
    source: str = "whale_alert"
    
    def __post_init__(self):
        """Set timestamp if not provided."""
        if self.timestamp is None:
            self.timestamp = datetime.utcnow()
    
    @property
    def formatted_amount(self) -> str:
        """Format the amount with appropriate units."""
        if self.amount >= 1_000_000:
            return f"{self.amount / 1_000_000:.2f}M {self.token_symbol}"
        elif self.amount >= 1_000:
            return f"{self.amount / 1_000:.2f}K {self.token_symbol}"
        else:
            return f"{self.amount:.2f} {self.token_symbol}"
    
    @property
    def formatted_usd_amount(self) -> str:
        """Format the USD amount."""
        if self.amount_usd >= 1_000_000:
            return f"${self.amount_usd / 1_000_000:.2f}M"
        elif self.amount_usd >= 1_000:
            return f"${self.amount_usd / 1_000:.2f}K"
        else:
            return f"${self.amount_usd:.2f}"
    
    def get_severity_color(self) -> int:
        """Get color based on USD amount."""
        if self.amount_usd >= 10_000_000:  # $10M+
            return 0xFF0000  # Red
        elif self.amount_usd >= 1_000_000:  # $1M+
            return 0xFF6600  # Orange
        elif self.amount_usd >= 100_000:   # $100K+
            return 0xFFCC00  # Yellow
        else:
            return 0x00CCFF  # Light blue
    
    def to_embed(self) -> dict:
        """Convert the whale alert to a Discord embed."""
        title = f"🐋 {self.alert_type.value.replace('_', ' ').title()}"
        
        description = f"**{self.formatted_amount}** ({self.formatted_usd_amount})"
        
        if self.from_label and self.to_label:
            description += f"\n📤 **From:** {self.from_label}\n📥 **To:** {self.to_label}"
        elif self.from_label:
            description += f"\n📤 **From:** {self.from_label}"
        elif self.to_label:
            description += f"\n📥 **To:** {self.to_label}"
        
        embed = {
            "title": title,
            "description": description,
            "color": self.get_severity_color(),
            "timestamp": self.timestamp.isoformat() if self.timestamp else datetime.utcnow().isoformat(),
            "footer": {
                "text": f"Source: {self.source}"
            },
            "fields": []
        }
        
        if self.blockchain:
            embed["fields"].append({
                "name": "Blockchain",
                "value": self.blockchain.title(),
                "inline": True
            })
        
        if self.transaction_hash:
            short_hash = f"{self.transaction_hash[:8]}...{self.transaction_hash[-8:]}"
            embed["fields"].append({
                "name": "Transaction",
                "value": f"`{short_hash}`",
                "inline": True
            })
        
        return embed

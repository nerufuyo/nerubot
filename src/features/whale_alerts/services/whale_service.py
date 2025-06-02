"""Whale alerts service for monitoring large crypto transactions."""
import asyncio
import logging
import aiohttp
from datetime import datetime, timedelta
from typing import List, Optional, Dict, Any
from ..models.whale_alert import WhaleAlert, AlertType


class WhaleService:
    """Service to monitor whale alerts from various sources."""
    
    def __init__(self, whale_alert_api_key: Optional[str] = None):
        """Initialize whale service with API credentials."""
        self.logger = logging.getLogger(__name__)
        self.whale_alert_api_key = whale_alert_api_key
        self.base_url = "https://api.whale-alert.io/v1"
        
        # Minimum USD amounts for different alert types
        self.min_amounts = {
            AlertType.LARGE_TRANSFER: 100_000,       # $100K
            AlertType.DEX_SWAP: 50_000,              # $50K
            AlertType.EXCHANGE_DEPOSIT: 500_000,      # $500K
            AlertType.EXCHANGE_WITHDRAWAL: 500_000,   # $500K
            AlertType.WHALE_ACCUMULATION: 1_000_000,  # $1M
            AlertType.WHALE_DISTRIBUTION: 1_000_000,  # $1M
        }
        
        # Exchange labels for categorization
        self.known_exchanges = {
            'binance', 'coinbase', 'kraken', 'bitfinex', 'huobi', 'okex',
            'kucoin', 'bybit', 'ftx', 'crypto.com', 'gate.io', 'bittrex'
        }
        
        self.last_transaction_time = datetime.utcnow() - timedelta(hours=1)
        self.is_monitoring = False
        
    async def start_monitoring(self, callback=None):
        """Start monitoring whale alerts."""
        if not self.whale_alert_api_key:
            self.logger.warning("Whale Alert API key not provided, using mock data")
        
        self.is_monitoring = True
        self.logger.info("Starting whale alerts monitoring")
        
        while self.is_monitoring:
            try:
                new_alerts = await self.fetch_new_alerts()
                
                if new_alerts and callback:
                    await callback(new_alerts)
                    
                # Wait 5 minutes before next check
                await asyncio.sleep(300)
                
            except Exception as e:
                self.logger.error(f"Error in whale alerts monitoring: {e}")
                await asyncio.sleep(600)  # Wait 10 minutes on error
    
    def stop_monitoring(self):
        """Stop monitoring."""
        self.is_monitoring = False
        self.logger.info("Stopped whale alerts monitoring")
    
    async def fetch_new_alerts(self) -> List[WhaleAlert]:
        """Fetch new whale alerts from API or generate mock data."""
        if self.whale_alert_api_key:
            return await self._fetch_from_api()
        else:
            return await self._generate_mock_alerts()
    
    async def _fetch_from_api(self) -> List[WhaleAlert]:
        """Fetch alerts from Whale Alert API."""
        alerts = []
        
        try:
            # Calculate time range (last hour)
            start_time = int(self.last_transaction_time.timestamp())
            end_time = int(datetime.utcnow().timestamp())
            
            url = f"{self.base_url}/transactions"
            params = {
                'api_key': self.whale_alert_api_key,
                'start': start_time,
                'end': end_time,
                'min_value': min(self.min_amounts.values()),
                'limit': 100
            }
            
            async with aiohttp.ClientSession() as session:
                async with session.get(url, params=params) as response:
                    if response.status == 200:
                        data = await response.json()
                        transactions = data.get('transactions', [])
                        
                        for tx in transactions:
                            alert = await self._convert_transaction_to_alert(tx)
                            if alert:
                                alerts.append(alert)
                    else:
                        self.logger.error(f"Whale Alert API error: {response.status}")
            
            self.last_transaction_time = datetime.utcnow()
            
        except Exception as e:
            self.logger.error(f"Error fetching from Whale Alert API: {e}")
        
        return alerts
    
    async def _convert_transaction_to_alert(self, transaction: Dict[str, Any]) -> Optional[WhaleAlert]:
        """Convert API transaction to WhaleAlert object."""
        try:
            # Determine alert type based on transaction details
            from_label = transaction.get('from', {}).get('owner', '').lower()
            to_label = transaction.get('to', {}).get('owner', '').lower()
            
            alert_type = AlertType.LARGE_TRANSFER
            
            if any(exchange in from_label for exchange in self.known_exchanges):
                if any(exchange in to_label for exchange in self.known_exchanges):
                    alert_type = AlertType.LARGE_TRANSFER
                else:
                    alert_type = AlertType.EXCHANGE_WITHDRAWAL
            elif any(exchange in to_label for exchange in self.known_exchanges):
                alert_type = AlertType.EXCHANGE_DEPOSIT
            elif 'uniswap' in from_label or 'uniswap' in to_label:
                alert_type = AlertType.DEX_SWAP
            
            # Skip if amount is below threshold for this type
            amount_usd = transaction.get('amount_usd', 0)
            if amount_usd < self.min_amounts.get(alert_type, 0):
                return None
            
            return WhaleAlert(
                alert_id=transaction.get('id', ''),
                alert_type=alert_type,
                token_symbol=transaction.get('symbol', 'UNKNOWN'),
                token_name=transaction.get('token_name', 'Unknown Token'),
                amount=transaction.get('amount', 0),
                amount_usd=amount_usd,
                from_address=transaction.get('from', {}).get('address'),
                to_address=transaction.get('to', {}).get('address'),
                from_label=transaction.get('from', {}).get('owner_type'),
                to_label=transaction.get('to', {}).get('owner_type'),
                transaction_hash=transaction.get('hash'),
                blockchain=transaction.get('blockchain'),
                timestamp=datetime.fromtimestamp(transaction.get('timestamp', 0)),
                source="whale_alert_api"
            )
            
        except Exception as e:
            self.logger.error(f"Error converting transaction: {e}")
            return None
    
    async def _generate_mock_alerts(self) -> List[WhaleAlert]:
        """Generate mock whale alerts for testing."""
        import random
        
        alerts = []
        
        # Generate 0-3 random alerts
        num_alerts = random.randint(0, 3)
        
        tokens = [
            {"symbol": "BTC", "name": "Bitcoin"},
            {"symbol": "ETH", "name": "Ethereum"},
            {"symbol": "USDT", "name": "Tether"},
            {"symbol": "SOL", "name": "Solana"},
            {"symbol": "ADA", "name": "Cardano"},
        ]
        
        exchanges = ["Binance", "Coinbase", "Kraken", "Unknown Wallet"]
        blockchains = ["bitcoin", "ethereum", "solana", "polygon"]
        
        for i in range(num_alerts):
            token = random.choice(tokens)
            alert_type = random.choice(list(AlertType))
            
            # Generate realistic amounts based on token
            if token["symbol"] == "BTC":
                amount = random.uniform(10, 1000)
                usd_value = amount * random.uniform(40000, 70000)
            elif token["symbol"] == "ETH":
                amount = random.uniform(100, 10000)
                usd_value = amount * random.uniform(2000, 4000)
            else:
                amount = random.uniform(100000, 10000000)
                usd_value = amount * random.uniform(0.5, 2.0)
            
            # Ensure minimum thresholds are met
            if usd_value < self.min_amounts.get(alert_type, 100000):
                usd_value = self.min_amounts.get(alert_type, 100000) * random.uniform(1.1, 5.0)
            
            alert = WhaleAlert(
                alert_id=f"mock_{datetime.utcnow().timestamp()}_{i}",
                alert_type=alert_type,
                token_symbol=token["symbol"],
                token_name=token["name"],
                amount=amount,
                amount_usd=usd_value,
                from_label=random.choice(exchanges),
                to_label=random.choice(exchanges),
                transaction_hash=f"0x{''.join(random.choices('0123456789abcdef', k=64))}",
                blockchain=random.choice(blockchains),
                timestamp=datetime.utcnow() - timedelta(minutes=random.randint(1, 60)),
                source="mock_data"
            )
            
            alerts.append(alert)
        
        return alerts
    
    def get_statistics(self) -> Dict[str, Any]:
        """Get monitoring statistics."""
        return {
            "is_monitoring": self.is_monitoring,
            "last_check": self.last_transaction_time.isoformat(),
            "min_amounts": {k.value: v for k, v in self.min_amounts.items()},
            "api_enabled": bool(self.whale_alert_api_key)
        }
    
    def set_minimum_amount(self, alert_type: AlertType, amount: float):
        """Set minimum USD amount for an alert type."""
        self.min_amounts[alert_type] = amount
        self.logger.info(f"Set minimum amount for {alert_type.value} to ${amount:,.2f}")
    
    async def get_recent_transactions(self, limit: int = 10) -> List[WhaleAlert]:
        """Get recent whale transactions."""
        if self.whale_alert_api_key:
            return await self._fetch_recent_from_api(limit)
        else:
            return await self._generate_mock_alerts()
    
    async def _fetch_recent_from_api(self, limit: int) -> List[WhaleAlert]:
        """Fetch recent transactions from API."""
        alerts = []
        
        try:
            # Get last 6 hours of data
            start_time = int((datetime.utcnow() - timedelta(hours=6)).timestamp())
            end_time = int(datetime.utcnow().timestamp())
            
            url = f"{self.base_url}/transactions"
            params = {
                'api_key': self.whale_alert_api_key,
                'start': start_time,
                'end': end_time,
                'min_value': min(self.min_amounts.values()),
                'limit': limit
            }
            
            async with aiohttp.ClientSession() as session:
                async with session.get(url, params=params) as response:
                    if response.status == 200:
                        data = await response.json()
                        transactions = data.get('transactions', [])
                        
                        for tx in transactions:
                            alert = await self._convert_transaction_to_alert(tx)
                            if alert:
                                alerts.append(alert)
                    else:
                        self.logger.error(f"Whale Alert API error: {response.status}")
            
        except Exception as e:
            self.logger.error(f"Error fetching recent transactions: {e}")
        
        return alerts[:limit]

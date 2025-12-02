"""
Copyright (c) 2023-2025 RapidaAI
Author: Prashant Srivastav <prashant@rapida.ai>

Licensed under GPL-2.0 with Rapida Additional Terms.
See LICENSE.md for details or contact sales@rapida.ai for commercial use.
"""
import bridges.bridge_factory


def init_app(app):
    app.extensions['provider_client'] = bridges.bridge_factory.get_me_provider_client(app.config.get('PROVIDER_API_HOST', 'localhost:9002'))

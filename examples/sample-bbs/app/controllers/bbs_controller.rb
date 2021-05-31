require 'net/https'
require 'securerandom'
require 'openssl'
require 'base64'

class BbsController < ApplicationController
  def index; end

  def login
    uri = URI.parse("#{Settings.login[:server_addr]}/authapi/v1/project/#{Settings.login[:project]}/openid-connect/auth")

    state = SecureRandom.hex(12)
    verifier = SecureRandom.hex(128)
    challenge = get_code_challenge(verifier)
    redirect_uri = "#{request.remote_ip}/bbs/callback"

    queries = {
      'scope' => 'openid email',
      'response_type' => 'code',
      'client_id' => Settings.login[:client_id],
      'redirect_uri' => redirect_uri,
      'code_challenge' => challenge,
      'code_challenge_method' => 'S256',
      'state' => state
    }
    uri.query = URI.encode_www_form(queries)

    logger.debug("login redirect to #{uri}")
    # TODO
    # redirect to hekate server /auth with params
    redirect_to uri.to_s
  end

  def callback; end

  def show
    @messages = Message.all
  end

  def add
    Message.create(
      text: params[:text],
      userid: 0 # debug
    )

    redirect_to action: 'show'
  end

  private

  def get_code_challenge(verifier)
    # currentryl supported only S256
    digest = OpenSSL::Digest.new('sha256')
    Base64.urlsafe_encode64(digest.update(verifier).digest)
  end

  def exchange_code
    # TODO
    # exchange auth code to token
    # redirect to top page
  end

  def find_user_info
    # TODO
    # get user info from JWT in web storage or cookie
    # if not exists, redirect to index page
  end
end

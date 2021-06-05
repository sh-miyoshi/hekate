require 'net/https'
require 'securerandom'
require 'openssl'
require 'base64'
require 'json'
require 'jwt'

class BbsController < ApplicationController
  before_action :find_user_info, only: %i[show add]

  def index; end

  def login
    uri = URI.parse("#{Settings.login[:server_addr]}/authapi/v1/project/#{Settings.login[:project]}/openid-connect/auth")

    state = SecureRandom.hex(12)
    verifier = SecureRandom.hex(128)
    challenge = get_code_challenge(verifier)
    redirect_uri = "#{Settings.login[:bbs_addr]}/bbs/callback"

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

    session[:state] = state
    session[:verifier] = verifier

    logger.debug("login redirect to #{uri}")
    redirect_to uri.to_s
  end

  def callback
    code = request.query_parameters[:code]
    state = request.query_parameters[:state]
    exchange_code(code, state)
    redirect_to action: 'show'
  end

  def show
    @messages = Message.all
  end

  def add
    Message.create(
      text: params[:text],
      userid: @user_id
    )

    redirect_to action: 'show'
  end

  private

  def get_code_challenge(verifier)
    # currentryl supported only S256
    digest = OpenSSL::Digest.new('sha256')
    Base64.urlsafe_encode64(digest.update(verifier).digest).delete('=')
  end

  def token_request(params)
    uri = URI.parse("#{Settings.login[:server_addr]}/authapi/v1/project/#{Settings.login[:project]}/openid-connect/token")
    res = Net::HTTP.post_form(uri, params)
    logger.debug("code exchange response: #{res.body}")

    raise 'failed to got token' if res.code.to_i > 300

    JSON.parse(res.body)
  end

  def exchange_code(code, state)
    logger.debug("current state: #{session[:state]}, got state: #{state}")
    raise 'invalid authentication state got' if state != session[:state]

    params = {
      'grant_type' => 'authorization_code',
      'client_id' => Settings.login[:client_id],
      'code' => code,
      'code_verifier' => session[:verifier],
      'state' => state
    }

    info = token_request(params)
    now = Time.current
    token = JWT.decode(info['access_token'], nil, false)
    logger.debug("Access token info: #{token}")

    s = LoginSession.create(
      expires_at: now.since(info['expires_in'].to_i),
      user_name: token[0]['preferred_username'],
      user_id: token[0]['sub']
    )

    session[:id] = s.id
  end

  def find_user_info
    if session[:id].nil?
      logger.debug('no session')
      redirect_to action: 'index'
      return
    end

    s = LoginSession.find(session[:id])
    now = Time.current

    # redirect if refresh token is expired
    if now >= s.expires_at
      logger.debug("token was expired. now: #{now}, expired at: #{s.expires_at}")
      redirect_to action: 'index'
      return
    end

    @user_id = s.user_id
    @user_name = s.user_name
  end
end
